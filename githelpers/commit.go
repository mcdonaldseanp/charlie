package githelpers

import (
	"fmt"
	. "github.com/mcdonaldseanp/charlie/airer"
	. "github.com/mcdonaldseanp/charlie/utils"
)

func NewCommit() (*Airer) {
	airr := addAllToWorkTree()
	if airr != nil { return airr }
	// Use the shell 'git commit' so that it will open vi to edit the message
	return ExecAsShell("git", "commit")
}

func AddCommit(no_edit bool) (*Airer) {
	airr := addAllToWorkTree()
	if airr != nil { return airr }
	// Use the shell 'git commit' so that it will open vi to edit the message
	if no_edit {
		// This technically doesn't need to use asShell, but whatever
		return ExecAsShell("git", "commit", "--amend", "--no-edit")
	} else {
		return ExecAsShell("git", "commit", "--amend")
	}
}

func addAllToWorkTree() (*Airer) {
	wt, airr := OpenWorktree()
	if airr != nil { return airr }

	status, err := wt.Status()
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to check branch status!\n%s\n", err),
			err,
		}
	}
	if len(status) < 1 {
		return &Airer{
			ExecError,
			fmt.Sprintf("Working tree is clean, nothing to commit\n"),
			nil,
		}
	}
	for filename, _ := range status {
		_, err := wt.Add(filename)
		if err != nil {
			return &Airer{
				ExecError,
				fmt.Sprintf("Failed git operation with file %s!\n%s\n", filename, err),
				err,
			}
		}
	}
	return nil
}