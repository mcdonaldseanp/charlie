package githelpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/mcdonaldseanp/charlie/airer"
)

func OpenRepo() (*git.Repository, *airer.Airer) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to open PWD!\n%s", err),
			Origin:  err,
		}
	}
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to load git repo!\n%s", err),
			Origin:  err,
		}
	}
	return repo, nil
}

func OpenWorktree() (*git.Worktree, *airer.Airer) {
	repo, airr := OpenRepo()
	if airr != nil {
		return nil, airr
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to load work tree!\n%s", err),
			Origin:  err,
		}
	}
	// Load the global gitignore and ensure the excludes patterns
	// in the work tree include the global patterns
	global_patterns, err := gitignore.LoadGlobalPatterns(osfs.New(filepath.Dir("/")))
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to load work tree!\n%s", err),
			Origin:  err,
		}
	}
	wt.Excludes = global_patterns
	return wt, nil
}

func WorkTreeClean(wt *git.Worktree) (bool, *airer.Airer) {
	status, err := wt.Status()
	if err != nil {
		return false, &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to check branch status!\n%s", err),
			Origin:  err,
		}
	}
	return len(status) == 0, nil
}
