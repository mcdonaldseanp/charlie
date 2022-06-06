package githelpers

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/auth"
	"github.com/mcdonaldseanp/charlie/localexec"
	. "github.com/mcdonaldseanp/charlie/utils"
)

func SetBranch(branch_name string, clear bool, pull bool) *airer.Airer {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	airr := ValidateParams(
		[]Validator{
			Validator{"branch_name", branch_name, []ValidateType{NotEmpty}},
		})
	if airr != nil {
		return airr
	}
	wt, airr := OpenWorktree()
	if airr != nil {
		return airr
	}

	clean, airr := WorkTreeClean(wt)
	if airr != nil {
		return airr
	}
	if !clean && !clear {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Cannot switch branch when work tree is not clean"),
			Origin:  nil,
		}
	}
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch_name),
		Force:  clear,
		Keep:   false,
	})
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to check out branch!\n%s", err),
			Origin:  err,
		}
	}
	if pull {
		_, airr := auth.TryFixAuth("git", "pull")
		if airr != nil {
			return airr
		}
	}
	return nil
}

func GetPR(pr_name string, clear bool, git_remote string) *airer.Airer {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	airr := ValidateParams(
		[]Validator{
			Validator{"branch_name", pr_name, []ValidateType{NotEmpty, IsNumber}},
		})
	if airr != nil {
		return airr
	}
	new_branch_name := "PR" + pr_name
	airr = localexec.ExecAsShell("git", "fetch", git_remote, "pull/"+pr_name+"/head:"+new_branch_name)
	if airr != nil {
		return airr
	}
	return SetBranch(new_branch_name, clear, false)
}
