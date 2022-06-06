package validator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/mcdonaldseanp/charlie/airer"
)

// Validators should take a string that looks something like:
//
// {"name":"input_param","value":"some_value_to_validate","validate":["NotEmpty","IsNumber"]}

type Validator struct {
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Validate []string `json:"validate"`
}

func ValidateParams(params string) *airer.Airer {
	unmarshald_data := []Validator{}
	// YAML would have been more human readable, but the use of whitespace in
	// yaml as an actual separator makes it really ugly to create docstrings
	// with yaml in them (which is the primary way this is meant to be used:
	// you pass this function a json docstring that identifies what to
	// validate
	err := json.Unmarshal([]byte(params), &unmarshald_data)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to parse validator as yaml:\n%s", err),
			Origin:  err,
		}
	}
	for _, data := range unmarshald_data {
		for _, validate_type := range data.Validate {
			switch validate_type {
			case "NotEmpty":
				if !(len(data.Value) > 0) {
					return &airer.Airer{
						Kind:    airer.InvalidInput,
						Message: fmt.Sprintf("'%s' is empty", data.Name),
						Origin:  nil,
					}
				}
			case "IsNumber":
				matcher, _ := regexp.Compile(`^[\d]+$`)
				if !matcher.Match([]byte(data.Value)) {
					return &airer.Airer{
						Kind:    airer.InvalidInput,
						Message: fmt.Sprintf("'%s' is not a number, given %s", data.Name, data.Value),
						Origin:  nil,
					}
				}
			case "IsIP":
				matcher, _ := regexp.Compile(`^[\d\.]+$`)
				if !matcher.Match([]byte(data.Value)) {
					return &airer.Airer{
						Kind:    airer.InvalidInput,
						Message: fmt.Sprintf("'%s' is not an IP address, given %s", data.Name, data.Value),
						Origin:  nil,
					}
				}
			case "IsFile":
				files, err := filepath.Glob(data.Value)
				if err != nil {
					return &airer.Airer{
						Kind:    airer.InvalidInput,
						Message: fmt.Sprintf("Failed attempting to check if '%s' is a file or directory, failure:\n%s", data.Name, err),
						Origin:  nil,
					}
				}
				if len(files) < 1 {
					return &airer.Airer{
						Kind:    airer.InvalidInput,
						Message: fmt.Sprintf("'%s' is not a file or directory, given %s", data.Name, data.Value),
						Origin:  nil,
					}
				}
			default:
				return &airer.Airer{
					Kind:    airer.ExecError,
					Message: fmt.Sprintf("Unknown matcher: %s", validate_type),
					Origin:  nil,
				}
			}
		}
	}
	return nil
}
