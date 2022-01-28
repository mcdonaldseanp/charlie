package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/airer"
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/utils"
)


func SetBranch(branch_name string, clear bool, pull bool) (*airer.Airer) {
	wt, airr := utils.OpenWorktree()
	if airr != nil { return airr }

	clean, airr := utils.WorkTreeClean(wt)
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
		_, airr := auth.TryFixAuth("git", "pull")
		if airr != nil { return airr }
	}
	return nil
}
