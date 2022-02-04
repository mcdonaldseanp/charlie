package githelpers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/McdonaldSeanp/charlie/auth"
	. "github.com/McdonaldSeanp/charlie/utils"
  . "github.com/McdonaldSeanp/charlie/airer"
)

func SetBranch(branch_name string, clear bool, pull bool) (*Airer) {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	airr := ValidateParams(
		[]Validator {
			Validator{ "branch_name", branch_name, []ValidateType{ NotEmpty } },
		})
	if airr != nil { return airr }
	wt, airr := OpenWorktree()
	if airr != nil { return airr }

	clean, airr := WorkTreeClean(wt)
	if airr != nil { return airr }
	if !clean && !clear {
		return &Airer{
			ExecError,
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
		return &Airer{
			ExecError,
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

func GetPR(pr_name string, clear bool) (*Airer) {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	airr := ValidateParams(
		[]Validator {
			Validator{ "branch_name", pr_name, []ValidateType{ NotEmpty, IsNumber } },
		})
	if airr != nil { return airr }
	new_branch_name := "PR" + pr_name
	airr = ExecAsShell("git", "fetch", "upstream", "pull/" + pr_name + "/head:" + new_branch_name)
	if airr != nil { return airr }
	return SetBranch(new_branch_name, clear, false)
}
