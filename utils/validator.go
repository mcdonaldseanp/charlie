package utils

import (
	"regexp"
	"fmt"
	. "github.com/McdonaldSeanp/charlie/airer"
)

type ValidateType int

const (
	NotEmpty ValidateType = iota
	IsNumber
)

type Validator struct {
	Name string
	Value string
	Validate []ValidateType
}

func ValidateParams(params []Validator) (*Airer) {
	for _, data := range params {
		for _, validate_type := range data.Validate {
			switch validate_type {
				case NotEmpty:
					if !(len(data.Value) > 0) {
						return &Airer{
							InvalidInput,
							fmt.Sprintf("Parameter '%s' is empty", data.Name),
							nil,
						}
					}
				case IsNumber:
					matcher, _ := regexp.Compile(`\d+`)
					if !matcher.Match([]byte(data.Value)) {
						return &Airer{
							InvalidInput,
							fmt.Sprintf("Parameter '%s' is not a number, given %s", data.Name, data.Value),
							nil,
						}
					}
			}
		}
	}
	return nil
}