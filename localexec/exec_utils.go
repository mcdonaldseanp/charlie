package localexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/sanitize"
)

// ExecAsShell always writes everything to stderr so that
// any resulting functionality can return something useful
// to the CLI
func ExecAsShell(command_string string, args ...string) *airer.Airer {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		new_command, arr := findWSLPath(command_string)
		if arr != nil {
			return &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Could not convert windows path to wsl path: %s", arr),
				Origin:  arr,
			}
		}
		command_string = sanitize.ReplaceAllNewlines(new_command)
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	shell_command.Stdout = os.Stderr
	shell_command.Stderr = os.Stderr
	shell_command.Stdin = os.Stdin
	err := shell_command.Run()
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ShellError,
			Message: fmt.Sprintf("Command %s failed: %s\n", shell_command, err),
			Origin:  err,
		}
	}
	return nil
}

func ExecDetached(command_string string, args ...string) (*exec.Cmd, *airer.Airer) {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		new_command, arr := findWSLPath(command_string)
		if arr != nil {
			return nil, &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Could not convert windows path to wsl path: %s", arr),
				Origin:  arr,
			}
		}
		command_string = sanitize.ReplaceAllNewlines(new_command)
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	err := shell_command.Start()
	if err != nil {
		return nil, &airer.Airer{
			Kind:    airer.ShellError,
			Message: fmt.Sprintf("Command '%s' failed to start:\n%s", shell_command, err),
			Origin:  err,
		}
	}
	return shell_command, nil
}

func ExecReadOutput(command_string string, args ...string) (string, string, *airer.Airer) {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		new_command, arr := findWSLPath(command_string)
		if arr != nil {
			return "", "", &airer.Airer{
				Kind:    airer.ExecError,
				Message: fmt.Sprintf("Could not convert windows path to wsl path: %s", arr),
				Origin:  arr,
			}
		}
		command_string = sanitize.ReplaceAllNewlines(new_command)
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	shell_command.Stdout = &stdout
	shell_command.Stderr = &stderr
	err := shell_command.Run()
	output := stdout.String()
	logs := stderr.String()
	if err != nil {
		return output, logs, &airer.Airer{
			Kind:    airer.ShellError,
			Message: fmt.Sprintf("Command '%s' failed:\n%s\nstderr:\n%s", shell_command, err, logs),
			Origin:  err,
		}
	}
	return output, logs, nil
}

func findWSLPath(command_string string) (string, *airer.Airer) {
	wsl_path, _, arr := ExecReadOutput("wslpath", "-u", command_string)
	if arr != nil {
		return "", arr
	}
	return wsl_path, nil
}

func isWinPath(command_string string) bool {
	return strings.HasPrefix(command_string, "C:\\") || strings.HasPrefix(command_string, "C:/")
}
