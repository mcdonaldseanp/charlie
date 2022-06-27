package githelpers

import (
	"github.com/mcdonaldseanp/charlie/localexec"
)

func NewCommit(maybe_message string) error {
	arr := addAllCLI()
	if arr != nil {
		return arr
	}
	// Use the shell 'git commit' so that it will open vi to edit the message
	if len(maybe_message) > 0 {
		arr = localexec.ExecAsShell("git", "commit", "-m", maybe_message)
	} else {
		arr = localexec.ExecAsShell("git", "commit")
	}
	return arr
}

func AddCommit(no_edit bool) error {
	arr := addAllCLI()
	if arr != nil {
		return arr
	}
	// Use the shell 'git commit' so that it will open vi to edit the message
	if no_edit {
		// This technically doesn't need to use asShell, but whatever
		arr := localexec.ExecAsShell("git", "commit", "--amend", "--no-edit")
		return arr
	} else {
		arr := error(localexec.ExecAsShell("git", "commit", "--amend"))
		return arr
	}
}

func addAllCLI() error {
	return localexec.ExecAsShell("git", "add", "--all")
}
