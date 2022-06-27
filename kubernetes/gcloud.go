package kubernetes

import (
	"fmt"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/clibuild/validator"
)

type GKECluster string

func InitializeGcloud() *airer.Airer {
	airr := localexec.ExecAsShell("gcloud", "auth", "login", "--no-launch-browser")
	if airr != nil {
		return airr
	}
	airr = localexec.ExecAsShell("gcloud", "config", "set", "project", "engineering-scratchpad")
	return airr
}

func ResizeCluster(cluster_name string, num_nodes string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"cluster_name","value":"%s","validate":["NotEmpty"]},
			{"name":"num_nodes","value":"%s","validate":["NotEmpty", "IsNumber"]}
		]`,
		cluster_name,
		num_nodes,
	))
	if arr != nil {
		return arr
	}
	return localexec.ExecAsShell("gcloud", "container", "clusters", "resize", cluster_name, "--num-nodes", num_nodes)
}

func (gkec GKECluster) NewClusterOfType(conf_loc string, extra_flags []string) *airer.Airer {
	// Only cluster_name is required, so that's the only thing to validate
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"cluster_name","value":"%s","validate":["NotEmpty"]}
		]`,
		string(gkec),
	))
	if arr != nil {
		return arr
	}
	return localexec.ExecAsShell(
		"gcloud",
		"container",
		"clusters",
		"create",
		string(gkec),
		"--release-channel",
		"None",
		"--machine-type",
		"e2-custom-6-16384",
		"--num-nodes",
		"1",
		"--addons",
		"HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver",
		"--no-enable-autoupgrade",
		"--no-enable-autorepair",
		"--enable-network-policy",
	)
}

func (gkec GKECluster) RemoveClusterOfType() *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"cluster_name","value":"%s","validate":["NotEmpty"]}]`,
		string(gkec),
	))
	if arr != nil {
		return arr
	}
	return localexec.ExecAsShell("gcloud", "container", "clusters", "delete", string(gkec))
}
