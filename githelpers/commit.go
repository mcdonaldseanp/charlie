package githelpers

import (
	"fmt"
	"github.com/McdonaldSeanp/charlie/airer"
	"github.com/McdonaldSeanp/charlie/utils"
)

func CommitAll() (*airer.Airer) {
	airr := AddAllToWorkTree()
	if airr != nil { return airr }
	// Use the shell 'git commit' so that it will open vi to edit the message
	airr = utils.ExecAsShell("git", "commit")
	if airr != nil { return airr }
	return nil
}

func AddAllToWorkTree() (*airer.Airer) {
	wt, airr := OpenWorktree()
	if airr != nil { return airr }

	status, err := wt.Status()
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check branch status!\n%s\n", err),
			err,
		}
	}
	if len(status) < 1 {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Working tree is clean, nothing to commit\n"),
			nil,
		}
	}
	for filename, _ := range status {
		err := wt.AddGlob(filename)
		if err != nil {
			return &airer.Airer{
				airer.ExecError,
				fmt.Sprintf("Failed to 'git add' file %s!\n%s\n", filename, err),
				err,
			}
		}
	}
	return nil
}