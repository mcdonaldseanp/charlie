package githelpers

import (
	"os"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/McdonaldSeanp/charlie/airer"
)

func OpenRepo() (*git.Repository, *airer.Airer) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to open PWD!\n%s\n", err),
			err,
		}
	}
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load git repo!\n%s\n", err),
			err,
		}
	}
	return repo, nil
}

func OpenWorktree() (*git.Worktree, *airer.Airer) {
	repo, airr := OpenRepo()
	if airr != nil { return nil, airr }

	wt, err := repo.Worktree()
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s\n", err),
			err,
		}
	}
	return wt, nil
}