package utils

import (
  "regexp"
  "fmt"
  "path/filepath"
  . "github.com/McdonaldSeanp/kelly/airer"
)

type ValidateType int

const (
  NotEmpty ValidateType = iota
  IsNumber
  IsFile
  IsIP
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
          matcher, _ := regexp.Compile(`^[\d]+$`)
          if !matcher.Match([]byte(data.Value)) {
            return &Airer{
              InvalidInput,
              fmt.Sprintf("Parameter '%s' is not a number, given %s", data.Name, data.Value),
              nil,
            }
          }
        case IsIP:
          matcher, _ := regexp.Compile(`^[\d\.]+$`)
          if !matcher.Match([]byte(data.Value)) {
            return &Airer{
              InvalidInput,
              fmt.Sprintf("Parameter '%s' is not a number, given %s", data.Name, data.Value),
              nil,
            }
          }
        case IsFile:
          files, err := filepath.Glob(data.Value)
          if err != nil {
            return &Airer{
              InvalidInput,
              fmt.Sprintf("Failed attempting to check if '%s' is a file or directory, failure:\n%s", data.Name, err),
              nil,
            }
          }
          if len(files) < 1 {
            return &Airer{
              InvalidInput,
              fmt.Sprintf("Parameter '%s' is not a file or directory, given %s", data.Name, data.Value),
              nil,
            }
          }
      }
    }
  }
  return nil
}
