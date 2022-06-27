package main

import (
	"flag"
	"os"
	"strings"

	"github.com/mcdonaldseanp/charlie/auth"
	"github.com/mcdonaldseanp/charlie/container"
	"github.com/mcdonaldseanp/charlie/githelpers"
	"github.com/mcdonaldseanp/charlie/kubernetes"
	"github.com/mcdonaldseanp/charlie/version"
	"github.com/mcdonaldseanp/clibuild/cli"
)

func main() {
	// None of the commands below should call .Parse on any
	// flagsets directly. cli.ShouldHaveArgs() will call .Parse
	// on the flagset if it is passed one.
	//
	// Things need to be parsed inside cli.ShouldHaveArgs so that
	// the flag package can ignore any required commands
	// before parsing

	// Git branch shared flags
	git_branch_fs := flag.NewFlagSet("git_branch", flag.ExitOnError)
	clear_branch := git_branch_fs.Bool("clear", false, "delete any changes in work tree")

	// Git commit shared flags
	git_commit_fs := flag.NewFlagSet("git_commit", flag.ExitOnError)

	// Git commit shared flags
	yubikey_fs := flag.NewFlagSet("yubikey", flag.ExitOnError)
	yubikey_hw_id := yubikey_fs.String("hardware-id", os.Getenv("YUBIKEY_HW_ID"), "hardware id of the yubikey to attach, in the form VID:PID")

	// Kubernetes create shared flags
	k8s_create_fs := flag.NewFlagSet("k8s_create_fs", flag.ExitOnError)
	k8s_create_conf_loc := k8s_create_fs.String("conf-loc", "", "Location on disk where the configuration file is for the cluster create command")
	k8s_create_extra_flags := k8s_create_fs.String("extra-flags", "", "Any additional flags to be sent with the cluster creation command")

	// Container shared flags
	con_fs := flag.NewFlagSet("container", flag.ExitOnError)
	container_registry := con_fs.String("registry",
		os.Getenv("DEFAULT_CONTAINER_REGISTRY"),
		"set the container registry")

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	command_list := []cli.Command{
		{
			Verb:     "add",
			Noun:     "commit",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie add commit [FLAGS]"
				description := "Add all changes in the work tree to previous commit"
				no_edit := git_commit_fs.Bool("no-edit", false, "Commit all changes without changing commit message")
				cli.ShouldHaveArgs(0, usage, description, git_commit_fs)
				cli.HandleCommandError(
					githelpers.AddCommit(*no_edit),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb:     "connect",
			Noun:     "pod",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie connect pod [POD NAME] [PORT]"
				description := "Use kubectl port-forward to open connection to a k8s pod, PORT can be omitted\n if POD NAME is a cygnus service name"
				// If this ever gets passed a flagset, the if statement below
				// needs to check if os.Args[4] is expected to be a flag or not
				cli.ShouldHaveArgs(1, usage, description, nil)
				port_num := ""
				if len(os.Args) >= 5 {
					port_num = os.Args[4]
				}
				cli.HandleCommandError(
					kubernetes.ConnectPod(os.Args[3], port_num),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb:     "disconnect",
			Noun:     "pod",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie disconnect pod [POD NAME]"
				description := "stop port fowarding from a k8s pod"
				cli.ShouldHaveArgs(1, usage, description, nil)
				cli.HandleCommandError(
					kubernetes.DisconnectPod(os.Args[3]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb:     "dismount",
			Noun:     "yubikey",
			Supports: []string{"linux"},
			ExecutionFn: func() {
				usage := "charlie dismount yubikey"
				description := "Detach yubikey from WSL instance"
				cli.ShouldHaveArgs(0, usage, description, yubikey_fs)
				cli.HandleCommandError(
					auth.DismountYubikey(*yubikey_hw_id),
					usage,
					description,
					yubikey_fs,
				)
			},
		},
		{
			Verb:     "get",
			Noun:     "pr",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				git_remote := git_branch_fs.String("remote", "upstream", "git remote to pull PR from")
				usage := "charlie get pr [PR NUMBER] [FLAGS]"
				description := "Check out contents of a PR from github"
				cli.ShouldHaveArgs(1, usage, description, git_branch_fs)
				cli.HandleCommandError(
					githelpers.GetPR(os.Args[3], *clear_branch, *git_remote),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb:     "initialize",
			Noun:     "gcloud",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie initialize gcloud"
				description := "initialize and authorize gcloud CLI"
				cli.ShouldHaveArgs(0, usage, description, nil)
				cli.HandleCommandError(
					kubernetes.InitializeGcloud(),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb:     "mount",
			Noun:     "yubikey",
			Supports: []string{"linux"},
			ExecutionFn: func() {
				usage := "charlie mount yubikey"
				description := "Connect yubikey to WSL instance"
				cli.ShouldHaveArgs(0, usage, description, yubikey_fs)
				cli.HandleCommandError(
					auth.MountYubikey(*yubikey_hw_id),
					usage,
					description,
					yubikey_fs,
				)
			},
		},
		{
			Verb:     "new",
			Noun:     "commit",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie new commit [FLAGS]"
				description := "create new commit from all changes in the work tree"
				message := git_commit_fs.String("message", "", "Provide the commit message")
				cli.ShouldHaveArgs(0, usage, description, git_commit_fs)
				cli.HandleCommandError(
					githelpers.NewCommit(*message),
					usage,
					description,
					git_commit_fs,
				)
			},
		},
		{
			Verb:     "new",
			Noun:     "cluster",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie new cluster [TYPE] [NAME] [FLAGS]"
				description := "Create a new kubernetes cluster of the given TYPE with name NAME. \nTYPE should be one of 'gke','kind'"
				cli.ShouldHaveArgs(2, usage, description, k8s_create_fs)
				extra_flags := strings.Split(*k8s_create_extra_flags, ",")
				cli.HandleCommandError(
					kubernetes.NewCluster(os.Args[3], os.Args[4], *k8s_create_conf_loc, extra_flags),
					usage,
					description,
					k8s_create_fs,
				)
			},
		},
		{
			Verb:     "publish",
			Noun:     "container",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie publish container [CONTAINER NAME] [NEW TAG] [FLAGS]"
				description := "publish the container that was last built locally to a container registry.\nDefaults to using DEFAULT_CONTAINER_REGISTRY env var"
				cli.ShouldHaveArgs(2, usage, description, con_fs)
				cli.HandleCommandError(
					container.PublishContainer(os.Args[3], os.Args[4], *container_registry),
					usage,
					description,
					con_fs,
				)
			},
		},
		{
			Verb:     "remove",
			Noun:     "cluster",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie remove cluster [NAME]"
				description := "Remove GKE cluster. Defaults to removing cluster with name from MY_CLUSTER \nenv var"
				cli.ShouldHaveArgs(1, usage, description, nil)
				cli.HandleCommandError(
					kubernetes.RemoveCluster(os.Args[3]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb:     "resize",
			Noun:     "cluster",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie resize cluster [NAME] [SIZE]"
				description := "resize Kubernetes cluster NAME to given SIZE"
				cli.ShouldHaveArgs(1, usage, description, nil)
				cli.HandleCommandError(
					kubernetes.ResizeCluster(os.Args[3], os.Args[4]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb:     "set",
			Noun:     "branch",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				pull_branch := git_branch_fs.Bool("pull", false, "pull from upstream")
				usage := "charlie set branch [BRANCH NAME] [FLAGS]"
				description := "set git repo to new branch"
				cli.ShouldHaveArgs(1, usage, description, git_branch_fs)
				cli.HandleCommandError(
					githelpers.SetBranch(os.Args[3], *clear_branch, *pull_branch),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb:     "start",
			Noun:     "docker",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "charlie start docker"
				description := "start the docker service on localhost"
				cli.ShouldHaveArgs(0, usage, description, nil)
				cli.HandleCommandError(
					container.StartDocker(),
					usage,
					description,
					nil,
				)
			},
		},
	}

	cli.RunCommand("charlie", version.VERSION, command_list)
}
