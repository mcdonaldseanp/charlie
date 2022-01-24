package auth

import (
	"fmt"
	"strings"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func FindYubikeyBUSID() (string, *airer.Airer) {
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

func ConnectYubikey() (*airer.Airer) {
	airr := utils.StartService("usbipd")
	if airr != nil {
		return airr
	}
	bus_id, airr := FindYubikeyBUSID()
	if airr != nil {
		return airr
	}
	airr = utils.ExecAsShell("usbipd.exe", "wsl", "attach", "--busid", bus_id)
	if airr != nil { return airr }
	airr = utils.ExecAsShell("sudo", "service", "pcscd", "restart")
	if airr != nil { return airr }
	return nil
}