package container

import (
	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	. "github.com/mcdonaldseanp/charlie/utils"
	"github.com/mcdonaldseanp/charlie/winservice"
)

func StartDocker() *airer.Airer {
	airr := winservice.StartService("com.docker.service")
	if airr != nil {
		return airr
	}
	_, airr = localexec.ExecDetached("/c/Program Files/Docker/Docker/Docker Desktop.exe")
	return airr
}

func PublishContainer(name string, tag string, registry_url string) *airer.Airer {
	airr := ValidateParams(
		[]Validator{
			Validator{"name", name, []ValidateType{NotEmpty}},
			Validator{"tag", tag, []ValidateType{NotEmpty}},
			Validator{"registry_url", registry_url, []ValidateType{NotEmpty}},
		})
	if airr != nil {
		return airr
	}
	output, _, airr := localexec.ExecReadOutput("docker", "images", "-q")
	if airr != nil {
		return airr
	}
	last_image := FirstLine(output)
	full_tag := registry_url + "/" + name + ":" + tag
	airr = localexec.ExecAsShell("docker", "tag", last_image, full_tag)
	if airr != nil {
		return airr
	}
	airr = localexec.ExecAsShell("docker", "push", full_tag)
	return airr
}
