package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/airer"
)


func Setgitbranch(branch_name string, clear bool) (error) {
	wt, err := OpenWorktree()
	if err != nil { return err }

	clean, err := WorkTreeClean(wt)
	if err != nil { return err }
	if !clean && !clear {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Cannot switch branch when work tree is not clean"),
		}
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch_name),
		Force: clear,
		Keep: false,
	})
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check out branch!\n%s\n", err),
		}
	}
	return nil
}

func WorkTreeClean(wt *git.Worktree) (bool, error) {
	status, err := wt.Status()
	if err != nil {
		return false, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check branch status!\n%s\n", err),
		}
	}
	return len(status) == 0, nil
}