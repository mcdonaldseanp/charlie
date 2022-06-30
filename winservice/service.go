package winservice

import (
	"os/exec"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/lookout/localexec"
)

func StartService(service_name string) error {
	airr := localexec.ExecAsShell("net.exe", "start", service_name)
	if airr != nil {
		if exitError, ok := airr.(*errtype.ShellError).Origin.(*exec.ExitError); ok {
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
