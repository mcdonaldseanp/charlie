package cygnus

import (
	"fmt"
	"github.com/McdonaldSeanp/charlie/container"
	. "github.com/McdonaldSeanp/charlie/utils"
	. "github.com/McdonaldSeanp/charlie/airer"
)

func ConnectCygnusPod(podname string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "podname", podname, []ValidateType{ NotEmpty } },
		})
	if airr != nil { return airr }
	switch podname {
		case "director":
			return container.ConnectPod("pe-orchestration-services-0", "8143")
		default:
			return &Airer{
				ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}

func DisconnectCygnusPod(podname string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "podname", podname, []ValidateType{ NotEmpty } },
		})
	if airr != nil { return airr }
	switch podname {
		case "director":
			return container.DisconnectPod("pe-orchestration-services-0")
		default:
			return &Airer{
				ExecError,
				fmt.Sprintf("Unknown Cygnus pod!"),
				nil,
			}
	}
}
