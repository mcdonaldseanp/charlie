package airer

import (
	"fmt"
)

type AirerType int

const (
	ShellError AirerType = iota
	ExecError
)

func (ar AirerType) String() string {
	return []string{"Shell command failed:", "Execution failure:"}[ar]
}

type Airer struct {
	Kind AirerType
	Message string
	Origin error
}

func (e *Airer) Error() string {
	return fmt.Sprintf("%s %s\n", e.Kind, e.Message)
}