package auth

import (
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/find"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/winservice"
)

func findYubikeyBUSID() (string, *airer.Airer) {
	output, _, airr := localexec.ExecReadOutput("usbipd.exe", "wsl", "list")
	if airr != nil {
		return "", airr
	}
	substr := find.LineWithSubStr(output, "Smartcard Reader")
	if substr == "" {
		return "", &airer.Airer{
			Kind:    airer.ExecError,
			Message: "Unable to find Yubikey BUSID, cannot continue",
			Origin:  nil,
		}
	}
	// Double negative here: returns true if the line does not
	// contain "Not attached"
	if !strings.Contains(substr, "Not attached") {
		return "", &airer.Airer{
			Kind:    airer.CompletedError,
			Message: "Yubikey already attached",
			Origin:  nil,
		}
	}
	return strings.Split(substr, " ")[0], nil
}

func MountYubikey() *airer.Airer {
	airr := winservice.StartService("usbipd")
	if airr != nil {
		return airr
	}
	bus_id, airr := findYubikeyBUSID()
	if airr != nil {
		return airr
	}
	airr = localexec.ExecAsShell("usbipd.exe", "wsl", "attach", "--busid", bus_id)
	if airr != nil {
		return airr
	}
	return localexec.ExecAsShell("sudo", "service", "pcscd", "restart")
}

func TryFixAuth(attempt_command string, params ...string) (string, *airer.Airer) {
	output, _, airr := localexec.ExecReadOutput(attempt_command, params...)
	if airr != nil {
		airr = RepairYubikey()
		if airr != nil {
			return "", &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Attempted to repair yubikey connection but failed\n%s", airr.Message),
				Origin:  airr,
			}
		}
		// Make another attempt
		output, _, airr = localexec.ExecReadOutput(attempt_command, params...)
		if airr != nil {
			return "", &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Attempted to repair yubikey connection but failed\n%s", airr.Message),
				Origin:  airr,
			}
		}
	}
	return output, nil
}

func RepairYubikey() *airer.Airer {
	// Make an attempt to load the yubikey in case that's the problem
	airr := MountYubikey()
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
