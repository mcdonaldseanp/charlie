package utils

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/McdonaldSeanp/charlie/airer"
)

func ExecAsShell(shell_command *exec.Cmd) (error) {
	shell_command.Stdout = os.Stdout
	shell_command.Stderr = os.Stderr
	shell_command.Stdin = os.Stdin
	err := shell_command.Run()
	if err != nil {
		return &airer.Airer{
			airer.ShellError,
			fmt.Sprintf("Command failed: %s\n", err),
		}
	}
	return nil
}
