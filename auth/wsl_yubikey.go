package auth

import (
	"fmt"
	"bytes"
	"strings"
	"os/exec"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func FindYubikeyBUSID() (string, error) {
	cmd := exec.Command("usbipd.exe", "wsl", "list")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", &airer.Airer{
			airer.ShellError,
			fmt.Sprintf("ERROR '%s' failed: %s\n\nstderr: %s\n", cmd, err, string(stderr.Bytes())),
		}
	}
	outStrLines := string(stdout.Bytes())
	return strings.Split(utils.LineWithSubStr(outStrLines, "Smartcard Reader"), " ")[0], nil
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