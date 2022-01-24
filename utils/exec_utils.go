package utils

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"github.com/McdonaldSeanp/charlie/airer"
)

func ExecAsShell(shell_command *exec.Cmd) (*airer.Airer) {
	shell_command.Stdout = os.Stdout
	shell_command.Stderr = os.Stderr
	shell_command.Stdin = os.Stdin
	err := shell_command.Run()
	if err != nil {
		return &airer.Airer{
			airer.ShellError,
			fmt.Sprintf("Command %s failed: %s\n", shell_command, err),
			err,
		}
	}
	return nil
}

func ExecReadOutput(shell_command *exec.Cmd) (string, *airer.Airer) {
	var stdout, stderr bytes.Buffer
	shell_command.Stdout = &stdout
	shell_command.Stderr = &stderr
	err := shell_command.Run()
	output := string(stdout.Bytes())
	if err != nil {
		return output, &airer.Airer{
			airer.ShellError,
			fmt.Sprintf("ERROR '%s' failed: %s\n\nstderr: %s\n", shell_command, err, string(stderr.Bytes())),
			err,
		}
	}
	return output, nil
}

func ExecDetached(shell_command *exec.Cmd) (*airer.Airer) {
	err := shell_command.Start()
	if err != nil {
		return &airer.Airer{
			airer.ShellError,
			fmt.Sprintf("ERROR '%s' failed to start: %s\n", shell_command, err),
			err,
		}
	}
	return nil
}