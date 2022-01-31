package cygnus

import (
	"fmt"
	"github.com/McdonaldSeanp/charlie/container"
	"github.com/McdonaldSeanp/charlie/airer"
)

func ConnectCygnusPod(podname string) (*airer.Airer) {
	switch podname {
		case "director":
			return container.ConnectPod("pe-orchestration-services-0", "8143")
		default:
			return &airer.Airer{
				airer.ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}

func DisconnectCygnusPod(podname string) (*airer.Airer) {
	switch podname {
		case "director":
			return container.DisconnectPod("pe-orchestration-services-0")
		default:
			return &airer.Airer{
				airer.ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}
