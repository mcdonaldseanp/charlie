package container

import (
	"fmt"
	"os"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/find"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/validator"
	"github.com/mcdonaldseanp/charlie/winservice"
)

var WORKSPACE_BIND_MOUNT string = "/wsl/Workspace/"

func StartDocker() *airer.Airer {
	// Make sure that the workspace is bind mounted to the cross distro space so that
	// k8s things can mount from localhost
	//
	// If the location is already bound, do nothing.
	_, _, arr := localexec.ExecReadOutput("findmnt", "--mountpoint", WORKSPACE_BIND_MOUNT)
	// findmnt will fail if it doesn't find anything
	if arr != nil {
		fmt.Print("Mounting workspace\n")
		arr = localexec.ExecAsShell("sudo", "mkdir", "-p", WORKSPACE_BIND_MOUNT)
		if arr != nil {
			return arr
		}
		arr = localexec.ExecAsShell("sudo", "mount", "--bind", os.Getenv("HOME")+"/Workspace", WORKSPACE_BIND_MOUNT)
		if arr != nil {
			return arr
		}
	}
	arr = winservice.StartService("com.docker.service")
	if arr != nil {
		return arr
	}
	_, arr = localexec.ExecDetached("/c/Program Files/Docker/Docker/Docker Desktop.exe")
	return arr
}

func PublishContainer(name string, tag string, registry_url string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"name","value":"%s","validate":["NotEmpty"]},
			{"name":"tag","value":"%s","validate":["NotEmpty"]},
			{"name":"registry_url","value":"%s","validate":["NotEmpty"]}
		 ]`,
		name,
		tag,
		registry_url,
	))
	if arr != nil {
		return arr
	}

	output, _, airr := localexec.ExecReadOutput("docker", "images", "-q")
	if airr != nil {
		return airr
	}
	last_image := find.FirstLine(output)
	full_tag := registry_url + "/" + name + ":" + tag
	airr = localexec.ExecAsShell("docker", "tag", last_image, full_tag)
	if airr != nil {
		return airr
	}
	airr = localexec.ExecAsShell("docker", "push", full_tag)
	return airr
}
