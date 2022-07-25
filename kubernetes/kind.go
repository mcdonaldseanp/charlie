package kubernetes

import (
	"fmt"
	"os"

	"github.com/mcdonaldseanp/charlie/localfile"
	"github.com/mcdonaldseanp/charlie/replacers"
	"github.com/mcdonaldseanp/charlie/version"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/localexec"
	"github.com/mcdonaldseanp/lookout/remotedata"
)

type KindCluster string

func (kc KindCluster) NewClusterOfType(conf_loc string, extra_flags []string) error {
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"cluster_name","value":"%s","validate":["NotEmpty"]}
		]`,
		string(kc),
	))
	if err != nil {
		return err
	}

	if len(conf_loc) < 1 {
		raw_data, err := remotedata.Download(version.ReleaseArtifact("kind_config.yaml"))
		if err != nil {
			return err
		}
		data, err := replacers.ReplaceVarsWithEnv(raw_data)
		if err != nil {
			return err
		}
		tmpfile, arr := localfile.TempFile("kind_config.yaml", []byte(data))
		if arr != nil {
			return arr
		}
		defer os.Remove(tmpfile)
		conf_loc = tmpfile
	}
	return localexec.ExecAsShell("kind", "create", "cluster", "--config", conf_loc, "--name", string(kc))
}

func (kc KindCluster) RemoveClusterOfType() error {
	return localexec.ExecAsShell("kind", "delete", "cluster", "--name", string(kc))
}
