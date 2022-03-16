package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/McdonaldSeanp/kelly/auth"
	"github.com/McdonaldSeanp/kelly/container"
	"github.com/McdonaldSeanp/kelly/cygnus"
	"github.com/McdonaldSeanp/kelly/gcloud"
	"github.com/McdonaldSeanp/kelly/githelpers"
	. "github.com/McdonaldSeanp/kelly/airer"
)

type CLICommand struct {
	Verb string
	Noun string
	ExecutionFn func()
}

func shouldHaveArgs(num_args int, usage string, description string, flagset *flag.FlagSet) {
	real_args := num_args + 1
	passed_fs := flagset != nil
	for _, arg := range os.Args {
		if arg == "-h" {
			fmt.Fprintf(os.Stderr, "Usage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Fprintf(os.Stderr, "Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(0)
		}
	}
	if len(os.Args) < real_args {
		fmt.Fprintf(os.Stderr, "AIRER running command:\n\nInvalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
		if passed_fs {
			fmt.Fprintf(os.Stderr, "Available flags:\n")
			flagset.PrintDefaults()
		}
		os.Exit(1)
	} else if len(os.Args) > real_args && passed_fs {
		flagset.Parse(os.Args[real_args:])
	}
}

func handleCommandAirer(airr *Airer, usage string, description string, flagset *flag.FlagSet) {
	if airr != nil {
		if airr.Kind == InvalidInput {
			fmt.Fprintf(os.Stderr, "%s\nUsage:\n  %s\n\nDescription:\n  %s\n\n", airr, usage, description)
			if flagset != nil {
				flagset.PrintDefaults()
			}
		} else {
			fmt.Fprintf(os.Stderr, "AIRER running command:\n\n%s\n", airr)
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

	// Git branch shared flags
	git_branch_fs := flag.NewFlagSet("git_branch", flag.ExitOnError)
	clear_branch := git_branch_fs.Bool("clear", false, "delete any changes in work tree")

	// Git commit shared flags
	git_commit_fs := flag.NewFlagSet("git_commit", flag.ExitOnError)

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

  cygnus_fs := flag.NewFlagSet("cygnus", flag.ExitOnError)
	cy_cluster_name := cygnus_fs.String("cluster-name",
													  os.Getenv("MY_CLUSTER"),
														"Execute against a specific cluster")
	build_repo_loc := cygnus_fs.String("build-repo",
									os.Getenv("CYGNUS_BUILD_REPO_DIR"),
									"Location on disk of the repo where cygnus builds from")
	pull_latest := cygnus_fs.Bool("pull-latest",
									false,
									"Updates the build repo to HEAD of the main branch from upstream when set")


	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	command_list := []CLICommand{
		{ "add", "commit",
			func() {
				usage := "kelly new commit [FLAGS]"
				description := "Add all changes in the work tree to previous commit"
				no_edit := git_commit_fs.Bool("no-edit", false, "Commit all changes without changing commit message")
				shouldHaveArgs(2, usage, description, git_commit_fs)
				handleCommandAirer(
					githelpers.AddCommit(*no_edit),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{ "connect", "pod",
			func() {
				usage := "kelly connect pod [POD NAME] [PORT]"
				description := "Use kubectl port-forward to open connection to a k8s pod, PORT can be omitted\n if POD NAME is a cygnus service name"
				// If this ever gets passed a flagset, the if statement below
				// needs to check if os.Args[4] is expected to be a flag or not
				shouldHaveArgs(3, usage, description, nil)
				port_num := ""
				if len(os.Args) >= 5 {
					port_num = os.Args[4]
				}
				handleCommandAirer(
					container.ConnectPod(os.Args[3], port_num),
					usage,
					description,
					nil,
				)
			},
		},
		{ "deploy", "cygnus",
			func() {
				usage := "kelly deploy cygnus [FLAGS]"
				description := "Deploy local changes of Cygnus to GKE"
				shouldHaveArgs(2, usage, description, cygnus_fs)
				handleCommandAirer(
					cygnus.DeployCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
		{ "disconnect", "pod",
			func() {
				usage := "kelly disconnect pod [POD NAME]"
				description := "stop port fowarding from a k8s pod"
				shouldHaveArgs(3, usage, description, nil)
				handleCommandAirer(
					container.DisconnectPod(os.Args[3]),
					usage,
					description,
					nil,
				)
			},
		},
		{ "get", "pr",
			func() {
				usage := "kelly get pr [PR NUMBER] [FLAGS]"
				description := "Check out contents of a PR from github"
				shouldHaveArgs(3, usage, description, git_branch_fs)
				handleCommandAirer(
					githelpers.GetPR(os.Args[3], *clear_branch),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{ "initialize", "gcloud",
			func() {
				usage := "kelly initialize gcloud"
				description := "initialize and authorize gcloud CLI"
				shouldHaveArgs(2, usage, description, nil)
				handleCommandAirer(
					gcloud.InitializeGcloud(),
					usage,
					description,
					nil,
				)
			},
		},
		{ "install", "cygnus",
			func() {
				usage := "kelly install cygnus [FLAGS]"
				description := "Deploy a new instance of Cygnus to GKE"
				shouldHaveArgs(2, usage, description, cygnus_fs)
				handleCommandAirer(
					cygnus.InstallCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
		{ "mount", "yubikey",
			func() {
				usage := "kelly mount yubikey"
				description := "Connect yubikey to WSL instance"
				shouldHaveArgs(2, usage, description, nil)
				handleCommandAirer(
					auth.MountYubikey(),
					usage,
					description,
					nil,
				)
			},
		},
		{ "new", "commit",
			func() {
				usage := "kelly new commit"
				description := "create new commit from all changes in the work tree"
				shouldHaveArgs(2, usage, description, nil)
				handleCommandAirer(
					githelpers.NewCommit(),
					usage,
					description,
					nil,
				)
			},
		},
		{ "new", "cluster",
			func() {
				usage := "kelly new cluster [SIZE] [FLAGS]"
				description := "Create a new GKE cluster with the given SIZE of nodes. Defaults to creating\n cluster with name from MY_CLUSTER env var"
				shouldHaveArgs(3, usage, description, gcloud_fs)
				handleCommandAirer(
					gcloud.NewCluster(*cluster_name, os.Args[3]),
					usage,
					description,
					gcloud_fs,
				)
			},
		},
		{ "publish", "container",
			func() {
				usage := "kelly publish container [CONTAINER NAME] [NEW TAG] [FLAGS]"
				description := "publish the container that was last built locally to a container registry.\nDefaults to using DEFAULT_CONTAINER_REGISTRY env var"
				shouldHaveArgs(4, usage, description, con_fs)
				handleCommandAirer(
					container.PublishContainer(os.Args[3], os.Args[4], *container_registry),
					usage,
					description,
					con_fs,
				)
			},
		},
		{ "read", "kotsip",
			func() {
				usage := "kelly read kotsip"
				description := "Read the ip that KOTS_IP should be set to"
				shouldHaveArgs(2, usage, description, nil)
				output, airr := cygnus.ReadKOTSIP()
				if airr == nil {
					fmt.Printf("%s", output)
				}
				handleCommandAirer(
					airr,
					usage,
					description,
					nil,
				)
			},
		},
		{ "remove", "cluster",
			func() {
				usage := "kelly remove cluster [FLAGS]"
				description := "Remove GKE cluster. Defaults to removing cluster with name from MY_CLUSTER \nenv var"
				shouldHaveArgs(2, usage, description, gcloud_fs)
				handleCommandAirer(
					gcloud.RemoveCluster(*cluster_name),
					usage,
					description,
					gcloud_fs,
				)
			},
		},
		{ "repair", "yubikey",
			func() {
				usage := "kelly repair yubikey"
				description := "attempt to repair yubikey connection to WSL instance"
				shouldHaveArgs(2, usage, description, nil)
				handleCommandAirer(
					auth.RepairYubikey(),
					usage,
					description,
					nil,
				)
			},
		},
		{ "resize", "cluster",
			func() {
				usage := "kelly resize cluster [SIZE] [FLAGS]"
				description := "resize GKE cluster to given SIZE. Defaults to resizing cluster with name \nfrom MY_CLUSTER env var"
				shouldHaveArgs(3, usage, description, gcloud_fs)
				handleCommandAirer(
					gcloud.ResizeCluster(*cluster_name, os.Args[3]),
					usage,
					description,
					gcloud_fs,
				)
			},
		},
		{ "set", "branch",
			func() {
				pull_branch := git_branch_fs.Bool("pull", false, "pull from upstream")
				usage := "kelly set branch [BRANCH NAME] [FLAGS]"
				description := "set git repo to new branch"
				shouldHaveArgs(3, usage, description, git_branch_fs)
				handleCommandAirer(
					githelpers.SetBranch(os.Args[3], *clear_branch, *pull_branch),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{ "start", "docker",
			func() {
				usage := "kelly start docker"
				description := "start the docker service on localhost"
				shouldHaveArgs(2, usage, description, nil)
				handleCommandAirer(
					container.StartDocker(),
					usage,
					description,
					nil,
				)
			},
		},
		{ "uninstall", "cygnus",
			func() {
				usage := "kelly uninstall cygnus [FLAGS]"
				description := "Run destroy-application to tear down an existing cygnus instance"
				shouldHaveArgs(2, usage, description, cygnus_fs)
				handleCommandAirer(
					cygnus.UninstallCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
	}

	if len(os.Args) > 1 {
		for _, command := range command_list {
			if os.Args[1] == command.Verb && os.Args[2] == command.Noun {
				command.ExecutionFn()
			}
		}
	}
	fmt.Printf("Unknown command.\n\nUsage:\n  kelly [COMMAND] [OBJECT] [ARGUMENTS] [FLAGS]\n\n")
	fmt.Printf("Available commands:\n")
	for _, command := range command_list {
		fmt.Printf("    %s %s\n", command.Verb, command.Noun)
	}
	os.Exit(1)
}