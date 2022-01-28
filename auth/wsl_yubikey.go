package auth

import (
	"fmt"
	"strings"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func findYubikeyBUSID() (string, *airer.Airer) {
	output, airr := utils.ExecReadOutput("usbipd.exe", "wsl", "list")
	if airr != nil { return "", airr }
	substr := utils.LineWithSubStr(output, "Smartcard Reader")
	if substr == "" {
		return "", &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Unable to find Yubikey BUSID, cannot continue"),
			nil,
		}
	}
	// Double negative here: returns true if the line does not
	// contain "Not attached"
	if !strings.Contains(substr, "Not attached") {
		return "", &airer.Airer{
			airer.CompletedError,
			"Yubikey already attached",
			nil,
		}
	}
	return strings.Split(substr, " ")[0], nil
}

func MountYubikey() (*airer.Airer) {
	airr := utils.StartService("usbipd")
	if airr != nil {
		return airr
	}
	bus_id, airr := findYubikeyBUSID()
	if airr != nil {
		return airr
	}
	airr = utils.ExecAsShell("usbipd.exe", "wsl", "attach", "--busid", bus_id)
	if airr != nil { return airr }
	airr = utils.ExecAsShell("sudo", "service", "pcscd", "restart")
	if airr != nil { return airr }
	return nil
}

func TryFixAuth(attempt_command string, params ...string) (string, *airer.Airer) {
	output, airr := utils.ExecReadOutput(attempt_command, params...)
	if airr != nil {
		// Make an attempt to load the yubikey in case that's the problem
		inner_airr := MountYubikey()
		if inner_airr != nil {
			if inner_airr.Kind != airer.CompletedError {
				// Yubikey didn't load, can't pull
				return "", &airer.Airer{
					airer.ExecError,
					fmt.Sprintf("Command '%s' failed, attempted to fix auth but failed\n\nOriginal failure:\n%s\nYubikey failure:\n%s", attempt_command, airr.Message, inner_airr.Message),
					airr,
				}
			}
			// If it was already connected, we can continue
		}
		// Make another attempt
		output, inner_airr = utils.ExecReadOutput(attempt_command, params...)
		if inner_airr != nil {
			// One last thing to try
			_, super_inner_airr := utils.ExecReadOutput("gpg-connect-agent", "updatestartuptty", "/bye")
			if super_inner_airr != nil {
				// Definately fucked
				return "", &airer.Airer{
					airer.ExecError,
					fmt.Sprintf("Command '%s' failed, attempted to fix auth but failed\n\nOriginal failure:\n%s\nGPG failure:\n%s", attempt_command, inner_airr.Message, super_inner_airr.Message),
					inner_airr,
				}
			}
			// one. final. attempt.
			output, super_inner_airr = utils.ExecReadOutput(attempt_command, params...)
			if super_inner_airr != nil {
				// Definately fucked
				return "", super_inner_airr
			}
		}
	}
	return output, nil
}