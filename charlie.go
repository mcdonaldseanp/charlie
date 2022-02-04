package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/container"
	"github.com/McdonaldSeanp/charlie/cygnus"
	"github.com/McdonaldSeanp/charlie/gcloud"
	"github.com/McdonaldSeanp/charlie/githelpers"
	. "github.com/McdonaldSeanp/charlie/airer"
)
func shouldHaveArgs(num_args int, usage string, description string, flagset *flag.FlagSet) {
	real_args := num_args + 1
	passed_fs := flagset != nil
	for _, arg := range os.Args {
		if arg == "-h" {
			fmt.Printf("Usage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Printf("Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(0)
		}
	}
	if len(os.Args) < real_args {
		fmt.Printf("Invalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
		if passed_fs {
			fmt.Printf("Available flags:\n")
			flagset.PrintDefaults()
		}
		os.Exit(1)
	} else if len(os.Args) > real_args && passed_fs {
		flagset.Parse(os.Args[real_args:])
	}
}

func handleCommand(airr *Airer, usage string, description string, flagset *flag.FlagSet) {
	if airr != nil {
		if airr.Kind == InvalidInput {
			fmt.Printf("%s\nUsage:\n  %s\n\nDescription:\n  %s\n\n", airr, usage, description)
			if flagset != nil {
				flagset.PrintDefaults()
			}
		} else {
			fmt.Printf("Failed to run command:\n\n%s\n", airr)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	// None of the commands below should call .Parse on any
	// flagsets directly. shouldHaveArgs() will call .Parse
	// on the flagset if it is passed one.
	//
	// Things need to be parsed inside shouldHaveArgs so that
	// the flag package can ignore any required commands
	// before parsing

	// Git shared flags
	git_fs := flag.NewFlagSet("git", flag.ExitOnError)
	clear_branch := git_fs.Bool("clear", false, "delete any changes in work tree")

	// Gcloud shared flags
	gcloud_fs := flag.NewFlagSet("gcloud", flag.ExitOnError)
	cluster_name := gcloud_fs.String("cluster-name",
													  os.Getenv("MY_CLUSTER"),
														"Execute against a specific cluster")

	// Container shared flags
	con_fs := flag.NewFlagSet("container", flag.ExitOnError)
	container_registry := con_fs.String("registry",
																 os.Getenv("DEFAULT_CONTAINER_REGISTRY"),
																 "set the container registry")

  // Used in unknown command usage at the bottom. Set here so that
	// you don't forget to add it.
	command_list := map[string][]string{
		"connect": {"pod", "cypod"},
		"disconnect": {"pod", "cypod"},
		"get": {"pr"},
		"initialize": {"gcloud"},
		"mount": {"yubikey"},
		"new": {"commit"},
		"publish": {"container"},
		"repair": {"yubikey"},
		"resize": {"cluster"},
		"set": {"branch"},
		"start": {"docker"},
	}

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	if len(os.Args) > 1 {
		switch os.Args[1] {
			case "connect":
				switch os.Args[2] {
					case "pod":
						usage := "charlie connect pod [POD NAME] [PORT]"
						description := "Use kubectl port-forward to open connection to a k8s pod"
						shouldHaveArgs(4, usage, description, nil)
						handleCommand(
							container.ConnectPod(os.Args[3], os.Args[4]),
							usage,
							description,
							nil,
						)
					case "cypod":
						usage := "charlie connect cypod [CYGNUS SERVICE NAME]"
						description := "Use kubectl port-forward to open a connection to a Cygnus k8s pod"
						shouldHaveArgs(3, usage, description, nil)
						handleCommand(
							cygnus.ConnectCygnusPod(os.Args[3]),
							usage,
							description,
							nil,
						)
				}
			case "disconnect":
				switch os.Args[2] {
					case "pod":
						usage := "charlie disconnect pod [POD NAME]"
						description := "stop port fowarding from a k8s pod"
						shouldHaveArgs(3, usage, description, nil)
						handleCommand(
							container.DisconnectPod(os.Args[3]),
							usage,
							description,
							nil,
						)
					case "cypod":
						usage :=	"charlie disconnect cypod [CYGNUS SERVICE NAME]"
						description := "stop port forwarding from a Cygnus k8s pod"
						shouldHaveArgs(3, usage, description, nil)
						handleCommand(
							cygnus.DisconnectCygnusPod(os.Args[3]),
							usage,
							description,
							nil,
						)
				}
			case "get":
				switch os.Args[2] {
					case "pr":
						usage := "charlie get pr [PR NUMBER] [FLAGS]"
						description := "Check out contents of a PR from github"
						shouldHaveArgs(3, usage, description, git_fs)
						handleCommand(
							githelpers.GetPR(os.Args[3], *clear_branch),
							usage,
							description,
							git_fs,
						)
				}
			case "initialize":
				switch os.Args[2] {
					case "gcloud":
						// Don't need to check args, since this isn't passed anything
						handleCommand(
							gcloud.InitializeGcloud(),
							"charlie initialize gcloud",
							"initialize and authorize gcloud CLI",
							nil,
						)
				}
			case "mount":
				switch os.Args[2] {
					case "yubikey":
						// Don't need to check args, since this isn't passed anything
						handleCommand(
							auth.MountYubikey(),
							"charlie mount yubikey",
							"Connect yubikey to WSL instance",
							nil,
						)
				}
			case "new":
				switch os.Args[2] {
					case "commit":
						// Don't need to check args, since this isn't passed anything
						handleCommand(
							githelpers.NewCommit(),
							"charlie new commit",
							"create new commit from all changes in the work tree",
							nil,
						)
				}
			case "publish":
				switch os.Args[2] {
					case "container":
						usage := "charlie publish container [CONTAINER NAME] [NEW TAG] [FLAGS]"
						description := "publish the container that was last built locally to a container registry.\nDefaults to using DEFAULT_CONTAINER_REGISTRY env var"
						shouldHaveArgs(4, usage, description, con_fs)
						handleCommand(
							container.PublishContainer(os.Args[3], os.Args[4], *container_registry),
							usage,
							description,
							con_fs,
						)
				}
			case "repair":
				switch os.Args[2] {
					case "yubikey":
						// Don't need to check args, since this isn't passed anything
						handleCommand(
							auth.RepairYubikey(),
							"charlie repair yubikey",
							"attempt to repair yubikey connection to WSL instance",
							nil,
						)
				}
			case "resize":
				switch os.Args[2] {
					case "cluster":
						usage := "charlie resize cluster [SIZE] [FLAGS]"
						description := "resize GKE cluster to given SIZE. Defaults to resizing cluster with name from MY_CLUSTER env var"
						shouldHaveArgs(3, usage, description, gcloud_fs)
						handleCommand(
							gcloud.ResizeCluster(*cluster_name, os.Args[3]),
							usage,
							description,
							gcloud_fs,
						)
				}
			case "set":
				switch os.Args[2] {
					case "branch":
						pull_branch := git_fs.Bool("pull", false, "pull from upstream")
						usage := "charlie set branch [BRANCH NAME] [FLAGS]"
						description := "set git repo to new branch"
						shouldHaveArgs(3, usage, description, git_fs)
						handleCommand(
							githelpers.SetBranch(os.Args[3], *clear_branch, *pull_branch),
							usage,
							description,
							git_fs,
						)
				}
			case "start":
				switch os.Args[2] {
					case "docker":
						// Don't need to check args, since this isn't passed anything
						handleCommand(
							container.StartDocker(),
							"charlie start docker",
							"start the docker service on localhost",
							nil,
						)
				}
		}
	}
	fmt.Printf("Unknown command.\n\nUsage:\n  charlie [COMMAND] [OBJECT] [ARGUMENTS] [FLAGS]\n\n")
	fmt.Printf("Available commands:\n")
	for verb, nouns := range command_list {
		for _, noun := range nouns {
			fmt.Printf("    %s %s\n", verb, noun)
		}
	}
	os.Exit(1)
}