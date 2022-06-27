package auth

import (
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

func MountYubikey(hw_id string) error {
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

func DismountYubikey(hw_id string) error {
	if !yubikeyAttached(hw_id) {
		return &airer.Airer{
			Kind:    airer.CompletedError,
			Message: "Yubikey already detached",
			Origin:  nil,
		}
	}
	return localexec.ExecAsShell("usbipd.exe", "wsl", "detach", "--hardware-id", hw_id)
}
