package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/githelpers"
	"github.com/McdonaldSeanp/charlie/container"
)

func main() {
	// Things need to be parsed inside the switch so that
	// the flag package can ignore any required commands
	// before parsing
	fs := flag.NewFlagSet("cli", flag.ExitOnError)
	clear_branch := fs.Bool("clear", false, "use --clear with 'set branch' to delete any changes in work tree")

	switch os.Args[1] {
		case "commit":
			switch os.Args[2] {
				case "all":
					err := githelpers.CommitAll()
					if err != nil {
						fmt.Printf("Did not commit!\n%s\n", err)
						os.Exit(1)
					}
				default:
					fmt.Printf("Unknown noun!\n")
					os.Exit(1)
			}
		case "load":
			switch os.Args[2] {
				case "yubikey":
					err := auth.ConnectYubikey()
					if err != nil {
						fmt.Printf("Did not load yubikey!\n%s\n", err)
						os.Exit(1)
					}
				default:
					fmt.Printf("Unknown noun!\n")
					os.Exit(1)
			}
		case "set":
			switch os.Args[2] {
				case "branch":
					fs.Parse(os.Args[4:])
					err := githelpers.Setgitbranch(os.Args[3], *clear_branch)
					if err != nil {
						fmt.Printf("Did not set branch!\n%s\n", err)
						os.Exit(1)
					}
				default:
					fmt.Printf("Unknown noun!\n")
					os.Exit(1)
			}
		case "start":
			switch os.Args[2] {
			case "docker":
				err := container.StartDocker()
				if err != nil {
					fmt.Printf("Could not start docker!\n%s\n", err)
					os.Exit(1)
				}
			default:
				fmt.Printf("Unknown noun!\n")
				os.Exit(1)
		}
		default:
			fmt.Printf("Unknown command!\n")
			os.Exit(1)
	}
	os.Exit(0)
}