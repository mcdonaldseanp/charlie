package localexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localfile"
)

// ExecAsShell always writes everything to stderr so that
// any resulting functionality can return something useful
// to the CLI
func ExecAsShell(command_string string, args ...string) *airer.Airer {
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

func ExecReadOutput(executable string, args ...string) (string, string, *airer.Airer) {
	shell_command := exec.Command(executable, args...)
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

func ExecScriptReadOutput(executable string, script string, args []string) (string, string, *airer.Airer) {
	f, err := os.CreateTemp("", "regulator_script")
	if err != nil {
		return "", "", &airer.Airer{
			Kind:    airer.ShellError,
			Message: "Could not create tmp file!",
			Origin:  err,
		}
	}
	filename := f.Name()
	defer os.Remove(filename) // clean up
	localfile.OverwriteFile(filename, []byte(script))
	final_args := append([]string{filename}, args...)
	return ExecReadOutput(executable, final_args...)
}

func BuildAndRunCommand(executable string, file string, script string, args []string) (string, string, *airer.Airer) {
	var output, logs string
	var arr *airer.Airer
	if len(file) > 0 {
		final_args := append([]string{file}, args...)
		output, logs, arr = ExecReadOutput(executable, final_args...)
	} else if len(script) > 0 {
		output, logs, arr = ExecScriptReadOutput(executable, script, args)
	} else {
		output, logs, arr = ExecReadOutput(executable, args...)
	}
	if arr != nil {
		return output, logs, arr
	}

	return output, logs, nil
}
