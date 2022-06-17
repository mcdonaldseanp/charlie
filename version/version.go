package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localfile"
	"github.com/mcdonaldseanp/charlie/validator"
)

const VERSION string = "v0.0.12"

func readVersion(raw_bytes []byte) (string, *airer.Airer) {
	lines := strings.Split(string(raw_bytes), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "const VERSION string") {
			split_line := strings.Split(line, " ")
			ver := split_line[len(split_line)-1]
			ver = ver[1 : len(ver)-1]
			return ver, nil
		}
	}
	return "", &airer.Airer{
		Kind:    airer.ExecError,
		Message: "Could not find version",
		Origin:  nil,
	}
}

func nextZ(ver string) (string, *airer.Airer) {
	split_ver := strings.Split(ver, ".")
	z_release, err := strconv.Atoi(split_ver[2])
	if err != nil {
		return "", &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Could not bump to next Z version, atoi conversion failed", err),
			Origin:  err,
		}
	}
	z_release++
	split_ver[2] = strconv.Itoa(z_release)
	return strings.Join(split_ver, "."), nil
}

func UpdateVersion(version_file string, new_version string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"version_file","value":"%s","validate":["NotEmpty", "IsFile"]}
		 ]`,
		version_file,
	))
	if arr != nil {
		return arr
	}
	raw_bytes, arr := localfile.ReadFileInChunks(version_file)
	if arr != nil {
		return arr
	}
	if len(new_version) < 1 {
		old_version, arr := readVersion(raw_bytes)
		if arr != nil {
			return arr
		}
		new_version, arr = nextZ(old_version)
		if arr != nil {
			return arr
		}
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
	arr = localfile.OverwriteFile(version_file, []byte(result))
	if arr != nil {
		return arr
	}
	fmt.Print(new_version)
	return nil
}

func ReleaseArtifact(name string) string {
	return fmt.Sprintf("https://github.com/mcdonaldseanp/charlie/releases/download/%s/%s", VERSION, name)
}
