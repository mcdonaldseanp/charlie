package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/airer"
)


func Setgitbranch(branch_name string, clear bool, pull bool) (*airer.Airer) {
	wt, airr := OpenWorktree()
	if airr != nil { return airr }

	clean, airr := WorkTreeClean(wt)
	if airr != nil { return airr }
	if !clean && !clear {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Cannot switch branch when work tree is not clean"),
			nil,
		}
	}
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch_name),
		Force: clear,
		Keep: false,
	})
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check out branch!\n%s", err),
			err,
		}
	}
	if pull {
		_, airr := TryFixAuth("git", "pull")
		if airr != nil { return airr }
	}
	return nil
}

func WorkTreeClean(wt *git.Worktree) (bool, *airer.Airer) {
	status, err := wt.Status()
	if err != nil {
		return false, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check branch status!\n%s", err),
			err,
		}
	}
	return len(status) == 0, nil
}