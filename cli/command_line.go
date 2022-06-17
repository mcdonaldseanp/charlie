package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/version"
)

type Command struct {
	Verb        string
	Noun        string
	Supports    []string
	ExecutionFn func()
}

// shouldHaveArgs does two things:
// * validate that the number of args that aren't flags have been provided (i.e. the number of strings
//    after the command name that aren't flags)
// * parse the remaining flags
//
// If the wrong number of args is passed it prints helpful usage
func ShouldHaveArgs(num_args int, usage string, description string, flagset *flag.FlagSet) {
	real_args := num_args + 1
	passed_fs := flagset != nil
	for index, arg := range os.Args {
		if arg == "-h" {
			fmt.Fprintf(os.Stderr, "Usage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Fprintf(os.Stderr, "Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(0)
		}
		// None of the arguments required by the command should start with dashes, if they
		// do assume an arg is missing and this is a flag
		if index <= num_args && strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Error running command:\n\nInvalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Fprintf(os.Stderr, "Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(1)
		}
	}
	if len(os.Args) < real_args {
		fmt.Fprintf(os.Stderr, "Error running command:\n\nInvalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
		if passed_fs {
			fmt.Fprintf(os.Stderr, "Available flags:\n")
			flagset.PrintDefaults()
		}
		os.Exit(1)
	} else if len(os.Args) > real_args && passed_fs {
		flagset.Parse(os.Args[real_args:])
	}
}

// handleCommandAirer catches InvalidInput airer.Airers and prints usage
// if that was the error thrown. IF a different type of airer.Airer is thrown
// it just prints the error.
//
// If the command succeeds handleCommandAirer exits the whole go process
// with code 0
func HandleCommandAirer(arr *airer.Airer, usage string, description string, flagset *flag.FlagSet) {
	if arr != nil {
		switch arr.Kind {
		case airer.InvalidInput:
			fmt.Fprintf(os.Stderr, "%s\nUsage:\n  %s\n\nDescription:\n  %s\n\n", arr, usage, description)
			if flagset != nil {
				flagset.PrintDefaults()
			}
		case airer.CompletedError:
			fmt.Printf("%s\n", arr.Message)
			// Completed "errors" are treated as success to the tool
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error running command:\n\n%s\n", arr)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func osSupportsCommand(cmd Command) bool {
	for _, os_name := range cmd.Supports {
		if runtime.GOOS == os_name {
			return true
		}
	}
	return false
}

func printTopUsage(tool_name string, command_list []Command) {
	fmt.Printf("Usage:\n  %s [COMMAND] [OBJECT] [ARGUMENTS] [FLAGS]\n\nAvailable commands:\n", tool_name)
	for _, command := range command_list {
		if osSupportsCommand(command) {
			fmt.Printf("    %s %s\n", command.Verb, command.Noun)
		}
	}
}

func RunCommand(tool_name string, command_list []Command) {
	if len(os.Args) > 2 {
		for _, command := range command_list {
			if os.Args[1] == command.Verb && os.Args[2] == command.Noun && osSupportsCommand(command) {
				command.ExecutionFn()
			}
		}
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Fprintf(os.Stdout, "%s\n", version.VERSION)
			os.Exit(0)
		case "-h":
			printTopUsage(tool_name, command_list)
			os.Exit(0)
		}
	}

	// If we've arrived here, that means the args passed don't match an existing command
	// --version or -h
	fmt.Printf("Unknown %s command \"%s\"\n\n", tool_name, strings.Join(os.Args, " "))
	printTopUsage(tool_name, command_list)
	os.Exit(1)
}
