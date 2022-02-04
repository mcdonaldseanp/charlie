package gcloud

import (
	. "github.com/McdonaldSeanp/charlie/utils"
	. "github.com/McdonaldSeanp/charlie/airer"
)

func InitializeGcloud() (*Airer) {
	airr := ExecAsShell("gcloud", "auth", "login", "--no-launch-browser")
	if airr != nil { return airr }
	airr = ExecAsShell("gcloud", "config", "set", "project", "engineering-scratchpad")
	return airr
}

func ResizeCluster(cluster_name string, num_nodes string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "cluster_name", cluster_name, []ValidateType{ NotEmpty } },
			Validator{ "num_nodes", num_nodes, []ValidateType{ NotEmpty, IsNumber } },
		})
	if airr != nil { return airr }
	airr = ExecAsShell("gcloud", "container", "clusters", "resize", cluster_name, "--num-nodes", num_nodes)
	return airr
}

func NewCluster(cluster_name string, num_nodes string) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "cluster_name", cluster_name, []ValidateType{ NotEmpty } },
			Validator{ "num_nodes", num_nodes, []ValidateType{ NotEmpty, IsNumber } },
		})
	if airr != nil { return airr }
	return ExecAsShell(
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