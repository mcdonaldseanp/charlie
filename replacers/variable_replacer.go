package replacers

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mcdonaldseanp/charlie/airer"
)

func strippedVar(raw_str []byte) string {
	str := string(raw_str)
	for i := 0; i < 3; i++ {
		_, left := utf8.DecodeRuneInString(str)
		_, right := utf8.DecodeLastRuneInString(str)
		str = str[left : len(str)-right]
	}
	return str
}

func ReplaceVarsWithEnv(input []byte) (string, error) {
	var result string = string(input)
	var vars_seen []string = make([]string, 0)
	var missing_os_vars []string = make([]string, 0)

	// See if there are any vars to replace
	matcher, _ := regexp.Compile(`\$__[A-Za-z_\d]+__\$`)
	matches := matcher.FindAll(input, -1)

	for _, raw_match := range matches {
		env_var_name := strippedVar(raw_match)
		for _, old_var := range vars_seen {
			if env_var_name == old_var {
				continue
			}
		}
		value := os.Getenv(env_var_name)
		result = strings.Replace(result, string(raw_match), value, -1)
		if len(value) < 1 {
			missing_os_vars = append(missing_os_vars, env_var_name)
		}
		vars_seen = append(vars_seen, env_var_name)
	}
	if len(missing_os_vars) != 0 {
		return result, &airer.Airer{
			Kind:    airer.InvalidInput,
			Message: fmt.Sprintf("The following env vars were empty: %s", strings.Join(missing_os_vars, ", ")),
			Origin:  nil,
		}
	}
	return result, nil
}
