package githelpers

import (
	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
)

func NewCommit() *airer.Airer {
	airr := addAllCLI()
	if airr != nil {
		return airr
	}
	// Use the shell 'git commit' so that it will open vi to edit the message
	return localexec.ExecAsShell("git", "commit")
}

func AddCommit(no_edit bool) *airer.Airer {
	airr := addAllCLI()
	if airr != nil {
		return airr
	}
	// Use the shell 'git commit' so that it will open vi to edit the message
	if no_edit {
		// This technically doesn't need to use asShell, but whatever
		return localexec.ExecAsShell("git", "commit", "--amend", "--no-edit")
	} else {
		return localexec.ExecAsShell("git", "commit", "--amend")
	}
}

func addAllCLI() *airer.Airer {
	return localexec.ExecAsShell("git", "add", "--all")
}
