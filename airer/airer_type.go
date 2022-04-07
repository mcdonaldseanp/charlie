package airer

import (
	"fmt"
)

type AirerType int

const (
	ShellError AirerType = iota
	ExecError
	CompletedError
	InvalidInput
	RemoteExecError
)

func (ar AirerType) String() string {
	return []string{"Shell command failed:", "Execution failed:", "Already done:", "Invalid input:", "Remote execution failed:"}[ar]
}

// Your illiteracy has screwed us again, regulator!
//
// Airer is a custom error type that provides a
// Kind field for parsing different error types.
//
// The name is an intentional misspelling of
// the word error.
type Airer struct {
	Kind    AirerType
	Message string
	Origin  error
}

func (e *Airer) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("%s\n%s\n\nTrace:\n%s\n", e.Kind, e.Message, e.Origin)
	} else {
		return fmt.Sprintf("%s\n%s\n", e.Kind, e.Message)
	}
}
