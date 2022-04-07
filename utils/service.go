package utils

import (
	"os/exec"
	. "github.com/mcdonaldseanp/charlie/airer"
)

func StartService(service_name string) (*Airer) {
	airr := ExecAsShell("net.exe", "start", service_name)
	if airr != nil {
		if exitError, ok := airr.Origin.(*exec.ExitError); ok {
			// If the exit code was '2' then the service was already running
			if exitError.ExitCode() != 2 {
				return airr
			}
		} else {
			return airr
		}
	}
	return nil
}
