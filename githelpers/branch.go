package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/airer"
)


func Setgitbranch(branch_name string) (error) {
	wt, err := OpenWorktree()
	if err != nil { return err }

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch_name),
		Force: false,
	})
	if err != nil {
		return &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to check out branch!\n%s\n", err),
		}
	}
	return nil
}