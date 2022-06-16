package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mcdonaldseanp/charlie/auth"
	"github.com/mcdonaldseanp/charlie/cli"
	"github.com/mcdonaldseanp/charlie/container"
	"github.com/mcdonaldseanp/charlie/cygnus"
	"github.com/mcdonaldseanp/charlie/githelpers"
	"github.com/mcdonaldseanp/charlie/kubernetes"
	"github.com/mcdonaldseanp/charlie/version"
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

	version_fs := flag.NewFlagSet("version", flag.ExitOnError)
	new_version := version_fs.String("version", "", "The new version to set, should be of the form vX.X.X. Defaults to bumping the Z version one digit")

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	command_list := []cli.Command{
		{
			Verb: "add",
			Noun: "commit",
			ExecutionFn: func() {
				usage := "charlie new commit [FLAGS]"
				description := "Add all changes in the work tree to previous commit"
				no_edit := git_commit_fs.Bool("no-edit", false, "Commit all changes without changing commit message")
				cli.ShouldHaveArgs(2, usage, description, git_commit_fs)
				cli.HandleCommandAirer(
					githelpers.AddCommit(*no_edit),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb: "connect",
			Noun: "pod",
			ExecutionFn: func() {
				usage := "charlie connect pod [POD NAME] [PORT]"
				description := "Use kubectl port-forward to open connection to a k8s pod, PORT can be omitted\n if POD NAME is a cygnus service name"
				// If this ever gets passed a flagset, the if statement below
				// needs to check if os.Args[4] is expected to be a flag or not
				cli.ShouldHaveArgs(3, usage, description, nil)
				port_num := ""
				if len(os.Args) >= 5 {
					port_num = os.Args[4]
				}
				cli.HandleCommandAirer(
					kubernetes.ConnectPod(os.Args[3], port_num),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "deploy",
			Noun: "cygnus",
			ExecutionFn: func() {
				usage := "charlie deploy cygnus [FLAGS]"
				description := "Deploy local changes of Cygnus to GKE"
				cli.ShouldHaveArgs(2, usage, description, cygnus_fs)
				cli.HandleCommandAirer(
					cygnus.DeployCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
		{
			Verb: "disconnect",
			Noun: "pod",
			ExecutionFn: func() {
				usage := "charlie disconnect pod [POD NAME]"
				description := "stop port fowarding from a k8s pod"
				cli.ShouldHaveArgs(3, usage, description, nil)
				cli.HandleCommandAirer(
					kubernetes.DisconnectPod(os.Args[3]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "dismount",
			Noun: "yubikey",
			ExecutionFn: func() {
				usage := "charlie dismount yubikey"
				description := "Detach yubikey from WSL instance"
				cli.ShouldHaveArgs(2, usage, description, yubikey_fs)
				cli.HandleCommandAirer(
					auth.DismountYubikey(*yubikey_hw_id),
					usage,
					description,
					yubikey_fs,
				)
			},
		},
		{
			Verb: "get",
			Noun: "pr",
			ExecutionFn: func() {
				git_remote := git_branch_fs.String("remote", "upstream", "git remote to pull PR from")
				usage := "charlie get pr [PR NUMBER] [FLAGS]"
				description := "Check out contents of a PR from github"
				cli.ShouldHaveArgs(3, usage, description, git_branch_fs)
				cli.HandleCommandAirer(
					githelpers.GetPR(os.Args[3], *clear_branch, *git_remote),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb: "initialize",
			Noun: "gcloud",
			ExecutionFn: func() {
				usage := "charlie initialize gcloud"
				description := "initialize and authorize gcloud CLI"
				cli.ShouldHaveArgs(2, usage, description, nil)
				cli.HandleCommandAirer(
					kubernetes.InitializeGcloud(),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "install",
			Noun: "cygnus",
			ExecutionFn: func() {
				usage := "charlie install cygnus [FLAGS]"
				description := "Deploy a new instance of Cygnus to GKE"
				cli.ShouldHaveArgs(2, usage, description, cygnus_fs)
				cli.HandleCommandAirer(
					cygnus.InstallCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
		{
			Verb: "mount",
			Noun: "yubikey",
			ExecutionFn: func() {
				usage := "charlie mount yubikey"
				description := "Connect yubikey to WSL instance"
				cli.ShouldHaveArgs(2, usage, description, yubikey_fs)
				cli.HandleCommandAirer(
					auth.MountYubikey(*yubikey_hw_id),
					usage,
					description,
					yubikey_fs,
				)
			},
		},
		{
			Verb: "new",
			Noun: "commit",
			ExecutionFn: func() {
				usage := "charlie new commit [FLAGS]"
				description := "create new commit from all changes in the work tree"
				message := git_commit_fs.String("message", "", "Provide the commit message")
				cli.ShouldHaveArgs(2, usage, description, git_commit_fs)
				cli.HandleCommandAirer(
					githelpers.NewCommit(*message),
					usage,
					description,
					git_commit_fs,
				)
			},
		},
		{
			Verb: "new",
			Noun: "cluster",
			ExecutionFn: func() {
				usage := "charlie new cluster [TYPE] [NAME] [FLAGS]"
				description := "Create a new kubernetes cluster of the given TYPE with name NAME. \nTYPE should be one of 'gke','kind'"
				cli.ShouldHaveArgs(4, usage, description, k8s_create_fs)
				extra_flags := strings.Split(*k8s_create_extra_flags, ",")
				cli.HandleCommandAirer(
					kubernetes.NewCluster(os.Args[3], os.Args[4], *k8s_create_conf_loc, extra_flags),
					usage,
					description,
					k8s_create_fs,
				)
			},
		},
		{
			Verb: "publish",
			Noun: "container",
			ExecutionFn: func() {
				usage := "charlie publish container [CONTAINER NAME] [NEW TAG] [FLAGS]"
				description := "publish the container that was last built locally to a container registry.\nDefaults to using DEFAULT_CONTAINER_REGISTRY env var"
				cli.ShouldHaveArgs(4, usage, description, con_fs)
				cli.HandleCommandAirer(
					container.PublishContainer(os.Args[3], os.Args[4], *container_registry),
					usage,
					description,
					con_fs,
				)
			},
		},
		{
			Verb: "read",
			Noun: "kotsip",
			ExecutionFn: func() {
				usage := "charlie read kotsip"
				description := "Read the ip that KOTS_IP should be set to"
				cli.ShouldHaveArgs(2, usage, description, nil)
				output, airr := cygnus.ReadKOTSIP()
				if airr == nil {
					fmt.Printf("%s", output)
				}
				cli.HandleCommandAirer(
					airr,
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "remove",
			Noun: "cluster",
			ExecutionFn: func() {
				usage := "charlie remove cluster [NAME]"
				description := "Remove GKE cluster. Defaults to removing cluster with name from MY_CLUSTER \nenv var"
				cli.ShouldHaveArgs(3, usage, description, nil)
				cli.HandleCommandAirer(
					kubernetes.RemoveCluster(os.Args[3]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "repair",
			Noun: "yubikey",
			ExecutionFn: func() {
				usage := "charlie repair yubikey"
				description := "attempt to repair yubikey connection to WSL instance"
				cli.ShouldHaveArgs(2, usage, description, yubikey_fs)
				cli.HandleCommandAirer(
					auth.RepairYubikey(*yubikey_hw_id),
					usage,
					description,
					yubikey_fs,
				)
			},
		},
		{
			Verb: "resize",
			Noun: "cluster",
			ExecutionFn: func() {
				usage := "charlie resize cluster [NAME] [SIZE]"
				description := "resize Kubernetes cluster NAME to given SIZE"
				cli.ShouldHaveArgs(3, usage, description, nil)
				cli.HandleCommandAirer(
					kubernetes.ResizeCluster(os.Args[3], os.Args[4]),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "set",
			Noun: "branch",
			ExecutionFn: func() {
				pull_branch := git_branch_fs.Bool("pull", false, "pull from upstream")
				usage := "charlie set branch [BRANCH NAME] [FLAGS]"
				description := "set git repo to new branch"
				cli.ShouldHaveArgs(3, usage, description, git_branch_fs)
				cli.HandleCommandAirer(
					githelpers.SetBranch(os.Args[3], *clear_branch, *pull_branch),
					usage,
					description,
					git_branch_fs,
				)
			},
		},
		{
			Verb: "start",
			Noun: "docker",
			ExecutionFn: func() {
				usage := "charlie start docker"
				description := "start the docker service on localhost"
				cli.ShouldHaveArgs(2, usage, description, nil)
				cli.HandleCommandAirer(
					container.StartDocker(),
					usage,
					description,
					nil,
				)
			},
		},
		{
			Verb: "uninstall",
			Noun: "cygnus",
			ExecutionFn: func() {
				usage := "charlie uninstall cygnus [FLAGS]"
				description := "Run destroy-application to tear down an existing cygnus instance"
				cli.ShouldHaveArgs(2, usage, description, cygnus_fs)
				cli.HandleCommandAirer(
					cygnus.UninstallCygnus(*cy_cluster_name, *build_repo_loc, *pull_latest),
					usage,
					description,
					cygnus_fs,
				)
			},
		},
		{
			Verb: "update",
			Noun: "version",
			ExecutionFn: func() {
				usage := "charlie update version [VERSION FILE] [FLAGS]"
				description := "Update charlie's version"
				cli.ShouldHaveArgs(3, usage, description, version_fs)
				cli.HandleCommandAirer(
					version.UpdateVersion(os.Args[3], *new_version),
					usage,
					description,
					version_fs,
				)
			},
		},
	}

	cli.RunCommand("charlie", command_list)
}
