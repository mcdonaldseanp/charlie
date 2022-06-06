package find

import (
	"strings"
)

func LineWithSubStr(str string, substr string) string {
	lines := strings.Split(str, "\n")
	var result string
	for _, line := range lines {
		if strings.Contains(line, substr) {
			result = line
			break
		}
	}
	return result
}

func FirstLine(str string) string {
	return strings.Split(str, "\n")[0]
}
