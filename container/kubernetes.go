package container

import (
	"strconv"
	"os"
	"fmt"
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func forwardPodPort(podname string, port string) (*airer.Airer) {
	forwarded_pods_file := os.Getenv("HOME") + "/.forwarded_pods"
	cmd, airr := utils.ExecDetached("kubectl", "port-forward", "pod/" + podname, port + ":" + port)
	if airr != nil { return airr }
	pid := strconv.Itoa(cmd.Process.Pid)
	forwarded_pods, airr := utils.ReadJSONFile(forwarded_pods_file)
	if airr != nil { return airr }
	if forwarded_pods == nil {
		forwarded_pods = make(map[string]interface {})
	}
	forwarded_pods[podname] = pid
	return utils.OverwriteJSONFile(forwarded_pods_file, forwarded_pods)
}

func stopForwardingPod(podname string) (*airer.Airer) {
	forwarded_pods_file := os.Getenv("HOME") + "/.forwarded_pods"
	forwarded_pods, airr := utils.ReadJSONFile(forwarded_pods_file)
	if airr != nil { return airr }
	if forwarded_pods[podname] == nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Nothing forwarded for pod %s", podname),
			nil,
		}
	}
	pid, err := strconv.Atoi(fmt.Sprintf("%v", forwarded_pods[podname]))
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed atoi conversion for %s", forwarded_pods[podname]),
			err,
		}
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to find process with PID %s", forwarded_pods[podname]),
			err,
		}
	}
	err = proc.Kill()
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to kill process with PID %s", forwarded_pods[podname]),
			err,
		}
	}
	forwarded_pods[podname] = nil
	return utils.OverwriteJSONFile(forwarded_pods_file, forwarded_pods)
}

func ConnectPod(podname string) (*airer.Airer) {
	switch podname {
		case "director":
			return forwardPodPort("pe-orchestration-services-0", "8143")
		default:
			return &airer.Airer{
				airer.ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}

func DisconnectPod(podname string) (*airer.Airer) {
	switch podname {
		case "director":
			return stopForwardingPod("pe-orchestration-services-0")
		default:
			return &airer.Airer{
				airer.ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}
