package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/find"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/winservice"
)

func yubikeyAttached(hw_id string) bool {
	output, _, airr := localexec.ExecReadOutput("usbipd.exe", "wsl", "list")
	if airr != nil {
		return false
	}
	substr := find.LineWithSubStr(output, hw_id)
	if substr == "" {
		return false
	}
	// Double negative here: returns true if the line does not
	// contain "Not attached"
	if !strings.Contains(substr, "Not attached") {
		return true
	}
	return false
}

func MountYubikey(hw_id string) *airer.Airer {
	airr := winservice.StartService("usbipd")
	if airr != nil {
		return airr
	}
	if yubikeyAttached(hw_id) {
		return &airer.Airer{
			Kind:    airer.CompletedError,
			Message: "Yubikey already attached",
			Origin:  nil,
		}
	}
	airr = localexec.ExecAsShell("usbipd.exe", "wsl", "attach", "--hardware-id", hw_id)
	if airr != nil {
		return airr
	}
	// Sleep for a couple seconds before restarting pcscd so that usbipd has
	// a chance to load the key
	time.Sleep(2 * time.Second)
	return localexec.ExecAsShell("sudo", "service", "pcscd", "restart")
}

func DismountYubikey(hw_id string) *airer.Airer {
	if !yubikeyAttached(hw_id) {
		return &airer.Airer{
			Kind:    airer.CompletedError,
			Message: "Yubikey already detached",
			Origin:  nil,
		}
	}
	return localexec.ExecAsShell("usbipd.exe", "wsl", "detach", "--hardware-id", hw_id)
}

func RepairYubikey(hw_id string) *airer.Airer {
	// Make an attempt to load the yubikey in case that's the problem
	airr := MountYubikey(hw_id)
	if airr != nil {
		if airr.Kind != airer.CompletedError {
			// Yubikey didn't load, can't fix
			return &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Yubikey could not be mounted:\n%s", airr.Message),
				Origin:  airr,
			}
		}
		// If it was already connected, we can continue
	}
	// Also try to unfuck gpg
	_, _, airr = localexec.ExecReadOutput("gpg-connect-agent", "updatestartuptty", "/bye")
	if airr != nil {
		// Definately fucked
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Attempt to repair GPG failed:\n%s", airr.Message),
			Origin:  airr,
		}
	}
	return nil
}
