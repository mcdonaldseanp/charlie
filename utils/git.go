package utils

import (
	"os"
	"fmt"
	"path/filepath"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-billy/v5/osfs"
	. "github.com/mcdonaldseanp/charlie/airer"
)

func OpenRepo() (*git.Repository, *Airer) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, &Airer{
			ExecError,
			fmt.Sprintf("Failed to open PWD!\n%s", err),
			err,
		}
	}
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, &Airer{
			ExecError,
			fmt.Sprintf("Failed to load git repo!\n%s", err),
			err,
		}
	}
	return repo, nil
}

func OpenWorktree() (*git.Worktree, *Airer) {
	repo, airr := OpenRepo()
	if airr != nil { return nil, airr }

	wt, err := repo.Worktree()
	if err != nil {
		return nil, &Airer{
			ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s", err),
			err,
		}
	}
	// Load the global gitignore and ensure the excludes patterns
	// in the work tree include the global patterns
	global_patterns, err := gitignore.LoadGlobalPatterns(osfs.New(filepath.Dir("/")))
	if err != nil {
		return nil, &Airer{
			ExecError,
			fmt.Sprintf("Failed to load work tree!\n%s", err),
			err,
		}
	}
	wt.Excludes = global_patterns
	return wt, nil
}

func WorkTreeClean(wt *git.Worktree) (bool, *Airer) {
	status, err := wt.Status()
	if err != nil {
		return false, &Airer{
			ExecError,
			fmt.Sprintf("Failed to check branch status!\n%s", err),
			err,
		}
	}
	return len(status) == 0, nil
}