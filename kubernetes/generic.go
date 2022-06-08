package kubernetes

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/localfile"
	"github.com/mcdonaldseanp/charlie/validator"
)

type KubernetesClusterType interface {
	NewClusterOfType(string, []string) *airer.Airer
	RemoveClusterOfType() *airer.Airer
}

var K8S_CLUSTER_FILE string = os.Getenv("HOME") + "/.charlie/data/clusters.json"
var K8S_FORWARDED_PODS_FILE string = os.Getenv("HOME") + "/.charlie/data/forwarded_pods.json"

func chooseClusterType(cluster_name string, cluster_type string) (KubernetesClusterType, *airer.Airer) {
	switch cluster_type {
	case "gke":
		var cluster GKECluster = GKECluster(cluster_name)
		return cluster, nil
	case "kind":
		var cluster KindCluster = KindCluster(cluster_name)
		return cluster, nil
	default:
		return nil, &airer.Airer{
			Kind:    airer.InvalidInput,
			Message: fmt.Sprintf("Unknown cluster type '%s'", cluster_type),
			Origin:  nil,
		}
	}
}

func saveNewClusterData(cluster_name string, cluster_type string) *airer.Airer {
	var data map[string]string
	arr := localfile.ReadJSONFile(K8S_CLUSTER_FILE, &data)
	if arr != nil {
		return arr
	}
	if data == nil {
		data = make(map[string]string)
	}
	data[cluster_name] = cluster_type
	return localfile.OverwriteJSONFile(K8S_CLUSTER_FILE, &data)
}

func DeleteClusterData(cluster_name string) *airer.Airer {
	var data map[string]string
	arr := localfile.ReadJSONFile(K8S_CLUSTER_FILE, &data)
	if arr != nil {
		return arr
	}
	if data == nil {
		data = make(map[string]string)
	}
	delete(data, cluster_name)
	return localfile.OverwriteJSONFile(K8S_CLUSTER_FILE, &data)
}

func readClusterType(cluster_name string) (string, *airer.Airer) {
	var data map[string]string
	arr := localfile.ReadJSONFile(K8S_CLUSTER_FILE, &data)
	if arr != nil {
		return "", arr
	}
	this_cluster_type, ok := data[cluster_name]
	if ok == false {
		return "", &airer.Airer{
			Kind:    airer.InvalidInput,
			Message: fmt.Sprintf("cluster '%s' has no data stored, cannot identify type", cluster_name),
			Origin:  nil,
		}
	}
	return this_cluster_type, nil
}

func NewCluster(cluster_type string, cluster_name string, conf_loc string, extra_flags []string) *airer.Airer {
	cluster, arr := chooseClusterType(cluster_name, cluster_type)
	if arr != nil {
		return arr
	}
	arr = cluster.NewClusterOfType(conf_loc, extra_flags)
	if arr != nil {
		return arr
	}
	return saveNewClusterData(cluster_name, cluster_type)
}

func RemoveCluster(cluster_name string) *airer.Airer {
	cluster_type, arr := readClusterType(cluster_name)
	if arr != nil {
		return arr
	}
	cluster, arr := chooseClusterType(cluster_name, cluster_type)
	if arr != nil {
		return arr
	}
	arr = cluster.RemoveClusterOfType()
	if arr != nil {
		return arr
	}
	return DeleteClusterData(cluster_name)
}

func ConnectPod(podname string, port string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"podname","value":"%s","validate":["NotEmpty"]}]`,
		podname,
	))
	if arr != nil {
		return arr
	}

	// BE WARY OF INFINITE LOOPS DUE TO MATCHING NAMES
	//
	// If any case names in this switch match the podname param
	// passed to the recursive call this will loop infinitely
	switch podname {
	case "director":
		return ConnectPod("pe-orchestration-services-0", "8143")
	}
	// Validate port separately in case a cygnus name was passed and port
	// is empty on the first function call.
	arr = validator.ValidateParams(fmt.Sprintf(
		`[{"name":"port","value":"%s","validate":["NotEmpty","IsNumber"]}]`,
		port,
	))
	if arr != nil {
		return arr
	}
	var forwarded_pods map[string]string
	cmd, airr := localexec.ExecDetached("kubectl", "port-forward", "pod/"+podname, port+":"+port)
	if airr != nil {
		return airr
	}
	pid := strconv.Itoa(cmd.Process.Pid)
	airr = localfile.ReadJSONFile(K8S_FORWARDED_PODS_FILE, &forwarded_pods)
	if airr != nil {
		return airr
	}
	if forwarded_pods == nil {
		forwarded_pods = make(map[string]string)
	}
	forwarded_pods[podname] = pid
	return localfile.OverwriteJSONFile(K8S_FORWARDED_PODS_FILE, &forwarded_pods)
}

func DisconnectPod(podname string) *airer.Airer {
	arr := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"podname","value":"%s","validate":["NotEmpty"]}]`,
		podname,
	))
	if arr != nil {
		return arr
	}
	switch podname {
	case "director":
		podname = "pe-orchestration-services-0"
	}
	var forwarded_pods map[string]string
	airr := localfile.ReadJSONFile(K8S_FORWARDED_PODS_FILE, &forwarded_pods)
	if airr != nil {
		return airr
	}
	this_pod, ok := forwarded_pods[podname]
	if ok == false {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Nothing forwarded for pod %s", podname),
			Origin:  nil,
		}
	}
	pid, err := strconv.Atoi(fmt.Sprintf("%v", this_pod))
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed atoi conversion for %s", this_pod),
			Origin:  err,
		}
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to find process with PID %s", this_pod),
			Origin:  err,
		}
	}
	err = proc.Kill()
	if err != nil {
		return &airer.Airer{
			Kind:    airer.ExecError,
			Message: fmt.Sprintf("Failed to kill process with PID %s", this_pod),
			Origin:  err,
		}
	}
	delete(forwarded_pods, podname)
	return localfile.OverwriteJSONFile(K8S_FORWARDED_PODS_FILE, &forwarded_pods)
}
