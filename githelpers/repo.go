package githelpers

import (
	"os"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/McdonaldSeanp/charlie/airer"
)

func OpenRepo() (*git.Repository, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to open PWD!\n%s\n", err),
		}
	}
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load git repo!\n%s\n", err),
		}
	}
	return repo, nil
}

func OpenWorktree() (*git.Worktree, error) {
	repo, err := OpenRepo()
	if err != nil { return nil, err }

	wt, err := repo.Worktree()
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s\n", err),
		}
	}
	return wt, nil
}