package container

import (
	"strconv"
	"os"
	"fmt"
	. "github.com/McdonaldSeanp/charlie/utils"
	. "github.com/McdonaldSeanp/charlie/airer"
)

func ConnectPod(podname string, port string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "podname", podname, []ValidateType{ NotEmpty } },
		})
	if airr != nil { return airr }
	// BE WARY OF INFINITE LOOPS DUE TO MATCHING NAMES
	switch podname {
		case "director":
			return ConnectPod("pe-orchestration-services-0", "8143")
	}
	// Validate port separately in case a cygnus name was passed and port
	// is empty on the first function call.
	airr = ValidateParams(
		[]Validator {
			Validator{ "port", port, []ValidateType{ NotEmpty, IsNumber } },
		})
	if airr != nil { return airr }
	forwarded_pods_file := os.Getenv("HOME") + "/.forwarded_pods"
	cmd, airr := ExecDetached("kubectl", "port-forward", "pod/" + podname, port + ":" + port)
	if airr != nil { return airr }
	pid := strconv.Itoa(cmd.Process.Pid)
	forwarded_pods, airr := ReadJSONFile(forwarded_pods_file)
	if airr != nil { return airr }
	if forwarded_pods == nil {
		forwarded_pods = make(map[string]interface {})
	}
	forwarded_pods[podname] = pid
	return OverwriteJSONFile(forwarded_pods_file, forwarded_pods)
}

func DisconnectPod(podname string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "podname", podname, []ValidateType{ NotEmpty } },
		})
	if airr != nil { return airr }
	switch podname {
		case "director":
			podname = "pe-orchestration-services-0"
	}
	forwarded_pods_file := os.Getenv("HOME") + "/.forwarded_pods"
	forwarded_pods, airr := ReadJSONFile(forwarded_pods_file)
	if airr != nil { return airr }
	if forwarded_pods[podname] == nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Nothing forwarded for pod %s", podname),
			nil,
		}
	}
	pid, err := strconv.Atoi(fmt.Sprintf("%v", forwarded_pods[podname]))
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed atoi conversion for %s", forwarded_pods[podname]),
			err,
		}
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to find process with PID %s", forwarded_pods[podname]),
			err,
		}
	}
	err = proc.Kill()
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to kill process with PID %s", forwarded_pods[podname]),
			err,
		}
	}
	forwarded_pods[podname] = nil
	return OverwriteJSONFile(forwarded_pods_file, forwarded_pods)
}
