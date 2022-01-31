package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/githelpers"
	"github.com/McdonaldSeanp/charlie/container"
	"github.com/McdonaldSeanp/charlie/cygnus"
)

func main() {
	// Things need to be parsed inside the switch so that
	// the flag package can ignore any required commands
	// before parsing
	fs := flag.NewFlagSet("cli", flag.ExitOnError)
	clear_branch := fs.Bool("clear", false, "use --clear with 'set branch' to delete any changes in work tree")
	pull_branch := fs.Bool("pull", false, "use --pull with 'set branch' to pull from upstream")

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	switch os.Args[1] {
		case "connect":
			switch os.Args[2] {
				case "pod":
					err := container.ConnectPod(os.Args[3], os.Args[4])
					if err != nil {
						fmt.Printf("Did not forward pod!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
				case "cypod":
					err := cygnus.ConnectCygnusPod(os.Args[3])
					if err != nil {
						fmt.Printf("Did not forward pod!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "disconnect":
			switch os.Args[2] {
				case "pod":
					err := container.DisconnectPod(os.Args[3])
					if err != nil {
						fmt.Printf("Could not stop pod forwarding!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
				case "cypod":
					err := cygnus.DisconnectCygnusPod(os.Args[3])
					if err != nil {
						fmt.Printf("Could not stop pod forwarding!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "get":
			switch os.Args[2] {
				case "pr":
					fs.Parse(os.Args[4:])
					err := githelpers.GetPR(os.Args[3], *clear_branch)
					if err != nil {
						fmt.Printf("Did not get PR!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "initialize":
			switch os.Args[2] {
				case "gcloud":
					err := container.InitializeGcloud()
					if err != nil {
						fmt.Printf("Did not authorize!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "mount":
			switch os.Args[2] {
				case "yubikey":
					err := auth.MountYubikey()
					if err != nil {
						fmt.Printf("Did not mount yubikey!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "new":
			switch os.Args[2] {
				case "commit":
					err := githelpers.NewCommit()
					if err != nil {
						fmt.Printf("Did not commit!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "repair":
			switch os.Args[2] {
				case "yubikey":
					err := auth.RepairYubikey()
					if err != nil {
						fmt.Printf("Could not repair yubikey!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "resize":
			switch os.Args[2] {
				case "cluster":
					err := container.ResizeCluster(os.Getenv("MY_CLUSTER"), os.Args[3])
					if err != nil {
						fmt.Printf("Did not resize cluster!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "set":
			switch os.Args[2] {
				case "branch":
					fs.Parse(os.Args[4:])
					err := githelpers.SetBranch(os.Args[3], *clear_branch, *pull_branch)
					if err != nil {
						fmt.Printf("Did not set branch!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
		case "start":
			switch os.Args[2] {
				case "docker":
					err := container.StartDocker()
					if err != nil {
						fmt.Printf("Could not start docker!\n%s\n", err)
						os.Exit(1)
					}
					os.Exit(0)
			}
	}
	fmt.Printf("Unknown command!\n")
	os.Exit(1)
}