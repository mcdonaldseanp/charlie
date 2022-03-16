package auth

import (
	"fmt"
	"strings"
	. "github.com/McdonaldSeanp/kelly/utils"
	. "github.com/McdonaldSeanp/kelly/airer"
)

func findYubikeyBUSID() (string, *Airer) {
	output, airr := ExecReadOutput("usbipd.exe", "wsl", "list")
	if airr != nil { return "", airr }
	substr := LineWithSubStr(output, "Smartcard Reader")
	if substr == "" {
		return "", &Airer{
			ExecError,
			fmt.Sprintf("Unable to find Yubikey BUSID, cannot continue"),
			nil,
		}
	}
	// Double negative here: returns true if the line does not
	// contain "Not attached"
	if !strings.Contains(substr, "Not attached") {
		return "", &Airer{
			CompletedError,
			"Yubikey already attached",
			nil,
		}
	}
	return strings.Split(substr, " ")[0], nil
}

func MountYubikey() (*Airer) {
	airr := StartService("usbipd")
	if airr != nil {
		return airr
	}
	bus_id, airr := findYubikeyBUSID()
	if airr != nil {
		return airr
	}
	airr = ExecAsShell("usbipd.exe", "wsl", "attach", "--busid", bus_id)
	if airr != nil { return airr }
	airr = ExecAsShell("sudo", "service", "pcscd", "restart")
	return nil
}

func TryFixAuth(attempt_command string, params ...string) (string, *Airer) {
	output, airr := ExecReadOutput(attempt_command, params...)
	if airr != nil {
		airr = RepairYubikey()
		if airr != nil {
			return "", &Airer{
				ExecError,
				fmt.Sprintf("Attempted to repair yubikey connection but failed\n%s", airr.Message),
				airr,
			}
		}
		// Make another attempt
		output, airr = ExecReadOutput(attempt_command, params...)
		if airr != nil {
			return "", &Airer{
				ExecError,
				fmt.Sprintf("Attempted to repair yubikey connection but failed\n%s", airr.Message),
				airr,
			}
		}
	}
	return output, nil
}

func RepairYubikey() (*Airer) {
	// Make an attempt to load the yubikey in case that's the problem
	airr := MountYubikey()
	if airr != nil {
		if airr.Kind != CompletedError {
			// Yubikey didn't load, can't fix
			return &Airer{
				ExecError,
				fmt.Sprintf("Yubikey could not be mounted:\n%s", airr.Message),
				airr,
			}
		}
		// If it was already connected, we can continue
	}
	// Also try to unfuck gpg
	_, airr = ExecReadOutput("gpg-connect-agent", "updatestartuptty", "/bye")
	if airr != nil {
		// Definately fucked
		return &Airer{
			ExecError,
			fmt.Sprintf("Attempt to repair GPG failed:\n%s", airr.Message),
			airr,
		}
	}
	return nil
}