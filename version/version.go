package version

import (
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localfile"
	"github.com/mcdonaldseanp/charlie/validator"
)

const VERSION string = "v0.0.5"

func UpdateVersion(version_file string, new_version string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"version_file","value":"%s","validate":["NotEmpty", "IsFile"]},
			{"name":"new_version","value":"%s","validate":["NotEmpty"]}
		 ]`,
		version_file,
		new_version,
	))
	if arr != nil {
		return arr
	}
	raw_bytes, arr := localfile.ReadFileInChunks(version_file)
	if arr != nil {
		return arr
	}
	lines := strings.Split(string(raw_bytes), "\n")
	var result string
	for index, line := range lines {
		if strings.HasPrefix(line, "const VERSION string") {
			result = result + fmt.Sprintf("const VERSION string = \"%s\"\n", new_version)
		} else {
			// Don't allow any newlines toward the end of the file
			//
			// This avoids creating a new newline every time the command is run
			// and making more and more newlines the more the command is run
			if len(line) > 0 || index < len(lines)-2 {
				result = result + line + "\n"
			}
		}
	}
	return localfile.OverwriteFile(version_file, []byte(result))
}

func ReleaseArtifact(name string) string {
	return fmt.Sprintf("https://github.com/mcdonaldseanp/charlie/releases/download/%s/%s", VERSION, name)
}
