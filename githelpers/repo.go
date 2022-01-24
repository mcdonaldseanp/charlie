package githelpers

import (
	"os"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/McdonaldSeanp/charlie/airer"
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/utils"
)

func OpenRepo() (*git.Repository, *airer.Airer) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to open PWD!\n%s", err),
			err,
		}
	}
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return nil, &airer.Airer{
			airer.ExecError,
			fmt.Sprintf("Failed to load git repo!\n%s", err),
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
			fmt.Sprintf("Failed to load work tree!\n%s", err),
			err,
		}
	}
	return wt, nil
}

func TryFixAuth(attempt_command string, params ...string) (string, *airer.Airer) {
	output, airr := utils.ExecReadOutput(attempt_command, params...)
	if airr != nil {
		// Make an attempt to load the yubikey in case that's the problem
		inner_airr := auth.ConnectYubikey()
		if inner_airr != nil {
			if inner_airr.Kind != airer.CompletedError {
				// Yubikey didn't load, can't pull
				return "", &airer.Airer{
					airer.ExecError,
					fmt.Sprintf("Command '%s' failed, attempted to fix auth but failed\n\nOriginal failure:\n%s\nYubikey failure:\n%s", attempt_command, airr.Message, inner_airr.Message),
					airr,
				}
			}
			// If it was already connected, we can continue
		}
		// Make another attempt
		output, inner_airr = utils.ExecReadOutput(attempt_command, params...)
		if inner_airr != nil {
			// One last thing to try
			_, super_inner_airr := utils.ExecReadOutput("gpg-connect-agent", "updatestartuptty", "/bye")
			if super_inner_airr != nil {
				// Definately fucked
				return "", &airer.Airer{
					airer.ExecError,
					fmt.Sprintf("Command '%s' failed, attempted to fix auth but failed\n\nOriginal failure:\n%s\nGPG failure:\n%s", attempt_command, inner_airr.Message, super_inner_airr.Message),
					inner_airr,
				}
			}
			// one. final. attempt.
			output, super_inner_airr = utils.ExecReadOutput(attempt_command, params...)
			if super_inner_airr != nil {
				// Definately fucked
				return "", super_inner_airr
			}
		}
	}
	return output, nil
}