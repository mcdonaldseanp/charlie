package version

import (
	"fmt"
)

const VERSION string = "v0.0.14"

func ReleaseArtifact(name string) string {
	return fmt.Sprintf("https://github.com/mcdonaldseanp/charlie/releases/download/%s/%s", VERSION, name)
}
