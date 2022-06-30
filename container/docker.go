package container

import (
	"fmt"

	"github.com/mcdonaldseanp/charlie/find"
	"github.com/mcdonaldseanp/charlie/winservice"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/localexec"
)

func StartDocker() error {
	arr := winservice.StartService("com.docker.service")
	if arr != nil {
		return arr
	}
	_, arr = localexec.ExecDetached("C:/Program Files/Docker/Docker/Docker Desktop.exe")
	return arr
}

func PublishContainer(name string, tag string, registry_url string) error {
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
