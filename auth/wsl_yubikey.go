package auth

import (
	"fmt"
	"strings"
	"os/exec"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func FindYubikeyBUSID() (string, error) {
	output, err := utils.ExecReadOutput(exec.Command("usbipd.exe", "wsl", "list"))
	if err != nil { return "", err }
	return strings.Split(utils.LineWithSubStr(output, "Smartcard Reader"), " ")[0], nil
}

func ConnectYubikey() (error) {
	bus_id, err := FindYubikeyBUSID()
	if err != nil { return err }
	if bus_id == "" {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("ERROR Unable to find Yubikey BUSID, cannot continue\n"),
		}
	}
	err = utils.ExecAsShell(exec.Command("usbipd.exe", "wsl", "attach", "--busid", bus_id))
	if err != nil { return err }
	err = utils.ExecAsShell(exec.Command("sudo", "service", "pcscd", "restart"))
	if err != nil { return err }
	return nil
}