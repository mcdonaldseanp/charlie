package kubernetes

import (
	"fmt"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/validator"
	"github.com/mcdonaldseanp/charlie/version"
)

type KindCluster string

func (kc KindCluster) NewClusterOfType(conf_loc string, extra_flags []string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"cluster_name","value":"%s","validate":["NotEmpty"]}
		]`,
		string(kc),
	))
	if arr != nil {
		return arr
	}
	if len(conf_loc) < 1 {
		conf_loc = version.ReleaseArtifact("kind_config.yaml")
	}
	return localexec.ExecAsShell("kind", "create", "cluster", "--config", conf_loc, "--name", string(kc))
}

func (kc KindCluster) RemoveClusterOfType() *airer.Airer {
	return localexec.ExecAsShell("kind", "delete", "cluster", "--name", string(kc))
}
