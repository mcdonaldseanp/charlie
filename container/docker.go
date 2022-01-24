package container

import (
	"os/exec"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func StartDocker() (*airer.Airer) {
	airr := utils.StartService("com.docker.service")
	if airr != nil { return airr }
	airr = utils.ExecDetached(exec.Command("/c/Program Files/Docker/Docker/Docker Desktop.exe"))
	if airr != nil { return airr }
	return nil
}
