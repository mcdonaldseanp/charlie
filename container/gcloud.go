package container

import (
	"github.com/McdonaldSeanp/charlie/utils"
	"github.com/McdonaldSeanp/charlie/airer"
)

func AuthorizeGcloud() (*airer.Airer) {
	airr := utils.ExecAsShell("gcloud", "auth", "login", "--no-launch-browser")
	if airr != nil { return airr }
	airr = utils.ExecAsShell("gcloud", "config", "set", "project", "engineering-scratchpad")
	if airr != nil { return airr }
	return nil
}

func ResizeGKECluster(cluster_name string, nodes string) (*airer.Airer) {
	airr := utils.ExecAsShell("gcloud", "container", "clusters", "resize", cluster_name, "--num-nodes", nodes)
	if airr != nil { return airr }
	return nil
}