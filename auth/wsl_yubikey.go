package auth

import (
	"fmt"
	"strings"
	"os/exec"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func FindYubikeyBUSID() (string, *airer.Airer) {
	airr := utils.StartService("usbipd")
	if airr != nil {
		return "", airr
	}
	output, airr := utils.ExecReadOutput(exec.Command("usbipd.exe", "wsl", "list"))
	if airr != nil { return "", airr }
	return strings.Split(utils.LineWithSubStr(output, "Smartcard Reader"), " ")[0], nil
}

func ConnectYubikey() (*airer.Airer) {
	bus_id, airr := FindYubikeyBUSID()
	if airr != nil {
		return airr
	}
	if bus_id == "" {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("airrOR Unable to find Yubikey BUSID, cannot continue\n"),
			nil,
		}
	}
	airr = utils.ExecAsShell(exec.Command("usbipd.exe", "wsl", "attach", "--busid", bus_id))
	if airr != nil { return airr }
	airr = utils.ExecAsShell(exec.Command("sudo", "service", "pcscd", "restart"))
	if airr != nil { return airr }
	return nil
}