package gcloud

import (
	"fmt"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/validator"
)

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

func NewCluster(cluster_name string, num_nodes string) *airer.Airer {
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
	return localexec.ExecAsShell(
		"gcloud",
		"container",
		"clusters",
		"create",
		cluster_name,
		"--release-channel",
		"None",
		"--machine-type",
		"e2-custom-6-16384",
		"--num-nodes",
		num_nodes,
		"--addons",
		"HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver",
		"--no-enable-autoupgrade",
		"--no-enable-autorepair",
		"--enable-network-policy",
	)
}

func RemoveCluster(cluster_name string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"cluster_name","value":"%s","validate":["NotEmpty"]}]`,
		cluster_name,
	))
	if arr != nil {
		return arr
	}
	return localexec.ExecAsShell("gcloud", "container", "clusters", "delete", cluster_name)
}
