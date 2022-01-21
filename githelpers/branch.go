package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/airer"
)


func Setgitbranch(branch_name string) (error) {
	repo, err := OpenRepo()
	if err != nil { return err }

	wt, err := repo.Worktree()
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s\n", err),
		}
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch_name),
	})
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s\n", err),
		}
	}
	return nil
}