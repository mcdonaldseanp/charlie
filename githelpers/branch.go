package githelpers

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/clibuild/validator"
)

func SetBranch(branch_name string, clear bool, pull bool) *airer.Airer {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"branch_name","value":"%s","validate":["NotEmpty"]}]`,
		branch_name,
	))
	if arr != nil {
		return arr
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
			Message: "Cannot switch branch when work tree is not clean",
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
		arr := localexec.ExecAsShell("git", "pull")
		if arr != nil {
			return arr
		}
	}
	return nil
}

func GetPR(pr_name string, clear bool, git_remote string) *airer.Airer {
	// Don't need to validate bool params, there will be a type error
	// if anything other than bools are passed
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"pr_name","value":"%s","validate":["NotEmpty", "IsNumber"]}]`,
		pr_name,
	))
	if arr != nil {
		return arr
	}

	new_branch_name := "PR" + pr_name
	arr = localexec.ExecAsShell("git", "fetch", git_remote, "pull/"+pr_name+"/head:"+new_branch_name)
	if arr != nil {
		return arr
	}
	return SetBranch(new_branch_name, clear, false)
}
