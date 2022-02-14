package cygnus

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/McdonaldSeanp/charlie/githelpers"
	. "github.com/McdonaldSeanp/charlie/utils"
	. "github.com/McdonaldSeanp/charlie/airer"
)

func InstallCygnus(cluster_name string, build_repo_loc string, pull_latest bool) (*Airer) {
	original_dir, err := os.Getwd()
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to get cwd with err: %s", err),
			err,
		}
	}
	defer os.Chdir(original_dir)
	cygnusBuildContext(cluster_name, build_repo_loc, pull_latest, false)
	airr := ExecAsShell("make", "start-ingress-controller")
	if airr != nil { return airr }
	// KOTS_IP has to be set outside the BuildContext helper since
	// on initial installation it won't exist and we need to run
	// start-ingress-controller before fetching the IP but _after_
	// the rest of the buildContext helper functionality
	err = os.Setenv("KOTS_IP", fetchKOTSIP())
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to set env var KOTS_IP with err: %s", err),
			err,
		}
	}
	fmt.Fprintf(os.Stderr, "KOTS IP IS NOW: %s\n", os.Getenv("KOTS_IP"))
	return ExecAsShell("make", "apply")
}

func UninstallCygnus(cluster_name string, build_repo_loc string, pull_latest bool) (*Airer) {
	original_dir, err := os.Getwd()
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to get cwd with err: %s", err),
			err,
		}
	}
	defer os.Chdir(original_dir)
	cygnusBuildContext(cluster_name, build_repo_loc, pull_latest, true)
	return ExecAsShell("make", "destroy-application")
}

func DeployCygnus(cluster_name string, build_repo_loc string, pull_latest bool) (*Airer) {
	original_dir, err := os.Getwd()
	if err != nil {
		return &Airer{
			ExecError,
			fmt.Sprintf("Failed to get cwd with err: %s", err),
			err,
		}
	}
	defer os.Chdir(original_dir)
	cygnusBuildContext(cluster_name, build_repo_loc, pull_latest, true)
	return ExecAsShell("make", "apply")
}

// Validate the params and then sets up the processes context so that:
//      * MY_CLUSTER to set to input
// 			* KOTS_IP is set correctly if fetch_kots_ip is true
//      * cwd is the build repo location
//      * if 'pull latest' is true, reset the git branch to HEAD of upstream's main branch
func cygnusBuildContext(cluster_name string, build_repo_loc string, pull_latest bool, fetch_kots_ip bool) (*Airer) {
	airr := ValidateParams(
		[]Validator {
			Validator{ "cluster_name", cluster_name, []ValidateType{ NotEmpty } },
			Validator{ "build_repo_loc", build_repo_loc, []ValidateType{ NotEmpty, IsFile } },
		})
	if airr != nil { return airr }
	// Set MY_CLUSTER and KOTS_IP
	os.Setenv("MY_CLUSTER", cluster_name)
	if fetch_kots_ip {
		err := os.Setenv("KOTS_IP", fetchKOTSIP())
		if err != nil {
			return &Airer{
				ExecError,
				fmt.Sprintf("Failed to set env var KOTS_IP with err: %s", err),
				err,
			}
		}
		fmt.Fprintf(os.Stderr, "KOTS IP IS NOW: %s\n", os.Getenv("KOTS_IP"))
	}
	// Change location to build repo, then check out the main branch if "pull_latest"
	// was specified
	os.Chdir(build_repo_loc)
	if pull_latest {
		airr = githelpers.SetBranch("main", true, true)
		if airr != nil { return airr }
	}
	return nil
}

func fetchKOTSIP() string {
	fmt.Fprintf(os.Stderr, "Reading KOTS_IP.")
	for {
		output, _ := ExecReadOutput(
			"kubectl",
			"-n",
			"ingress-nginx",
			"get",
			"svc",
			"ingress-nginx-controller",
			"-o",
			"jsonpath='{.status.loadBalancer.ingress[0].ip}'",
		)
		// I could fix this if I cared about it :D
		output = strings.Trim(output, "'")
		airr := ValidateParams(
			[]Validator {
				Validator{ "kots_ip", output, []ValidateType{ NotEmpty, IsIP } },
			})
		if airr == nil {
			fmt.Fprintf(os.Stderr, "\n")
			return output
		}
		fmt.Fprintf(os.Stderr, ".")
		time.Sleep(2 * time.Second)
	}
}

func ReadKOTSIP() (*Airer) {
	fmt.Printf("%s", fetchKOTSIP())
	return nil
}