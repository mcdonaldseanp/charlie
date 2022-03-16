package utils

import (
  "fmt"
  "os"
  "os/exec"
  "bytes"
  . "github.com/McdonaldSeanp/kelly/airer"
)

// ExecAsShell always writes everything to stderr so that
// any resulting functionality can return something useful
// to the CLI
func ExecAsShell(command_string string, params ...string) (*Airer) {
  shell_command := exec.Command(command_string, params...)
  shell_command.Env = os.Environ()
  shell_command.Stdout = os.Stderr
  shell_command.Stderr = os.Stderr
  shell_command.Stdin = os.Stdin
  err := shell_command.Run()
  if err != nil {
    return &Airer{
      ShellError,
      fmt.Sprintf("Command %s failed: %s\n", shell_command, err),
      err,
    }
  }
  return nil
}

func ExecReadOutput(command_string string, params ...string) (string, *Airer) {
  shell_command := exec.Command(command_string, params...)
  shell_command.Env = os.Environ()
  var stdout, stderr bytes.Buffer
  shell_command.Stdout = &stdout
  shell_command.Stderr = &stderr
  err := shell_command.Run()
  output := string(stdout.Bytes())
  if err != nil {
    return output, &Airer{
      ShellError,
      fmt.Sprintf("Command '%s' failed:\n%s\nstderr:\n%s", shell_command, err, string(stderr.Bytes())),
      err,
    }
  }
  return output, nil
}

func ExecDetached(command_string string, params ...string) (*exec.Cmd, *Airer) {
  shell_command := exec.Command(command_string, params...)
  shell_command.Env = os.Environ()
  err := shell_command.Start()
  if err != nil {
    return nil, &Airer{
      ShellError,
      fmt.Sprintf("Command '%s' failed to start:\n%s", shell_command, err),
      err,
    }
  }
  return shell_command, nil
}
