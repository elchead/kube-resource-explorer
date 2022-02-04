package kube

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeresource "k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	clientset "k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/tools/clientcmd"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}

func GetClient(config *rest.Config) *KubeClient {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return NewKubeClient(clientset)
}

func GetNodesListAndMetrics(config *rest.Config) (*v1.NodeList, *versioned.Clientset) {
	k := GetClient(config)

	metricsclient, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := k.NodeList()
	if err != nil {
		panic(err.Error())
	}
	return nodes, metricsclient
}

// mem usage in Gi
type NodeStatus struct {
	name     string
	memReq   int64
	MemAlloc int64
}

func GetPodsTotalMemRequestsAndLimits(podList []corev1.Pod) (int64, int64) {
	reqs, limits := map[corev1.ResourceName]resource.Quantity{}, map[corev1.ResourceName]resource.Quantity{}
	for _, pod := range podList {
		podReqs, podLimits := resourcehelper.PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	memoryReqs, memoryLimits := reqs[corev1.ResourceMemory], limits[corev1.ResourceMemory]
	return memoryReqs.ScaledValue(resource.Scale(9)), memoryLimits.ScaledValue(resource.Scale(9))
}

func getPodsInChunks(c corev1client.PodInterface, initialOpts metav1.ListOptions) (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	err := runtimeresource.FollowContinue(&initialOpts,
		func(options metav1.ListOptions) (runtime.Object, error) {
			newList, err := c.List(context.TODO(), options)
			if err != nil {
				return nil, runtimeresource.EnhanceListError(err, options, corev1.ResourcePods.String())
			}
			podList.Items = append(podList.Items, newList.Items...)
			return newList, nil
		})
	return podList, err
}

func GetPodsByNode(d clientset.Clientset, name string, namespace string) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + name + ",status.phase!=" + string(corev1.PodSucceeded) + ",status.phase!=" + string(corev1.PodFailed))
	if err != nil {
		return nil, err
	}
	initialOpts := metav1.ListOptions{
		FieldSelector: fieldSelector.String(),
	}
	return getPodsInChunks(d.CoreV1().Pods(namespace), initialOpts)
}

// TODO check sort NOT WORKING FOR Memory
func GetNodesByUsage(nodes *v1.NodeList, metricsclient *versioned.Clientset) map[string]NodeStatus {
	res := map[string]NodeStatus{} //[]NodeStatus{}
	for _, node := range nodes.Items {
		mc, err := metricsclient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.GetName(), metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		// reqs, limits := getPodsTotalRequestsAndLimits(nodeNonTerminatedPodsList)
		// cpuReqs, cpuLimits, memoryReqs, memoryLimits, ephemeralstorageReqs, ephemeralstorageLimits :=
		// reqs[corev1.ResourceCPU], limits[corev1.ResourceCPU], reqs[corev1.ResourceMemory], limits[corev1.ResourceMemory], reqs[corev1.ResourceEphemeralStorage], limits[corev1.ResourceEphemeralStorage]
		usedMem := mc.Usage.Memory().ScaledValue(resource.Scale(9))
		allocatable := node.Status.Capacity
		if len(node.Status.Allocatable) > 0 {
			allocatable = node.Status.Allocatable
		}
		lim := allocatable.Memory().ScaledValue(resource.Scale(9))
		node := NodeStatus{node.GetName(), usedMem, lim}
		res[node.name] = node
		// res = append(res, node)
	}
	return res
}

func FilterNodesByUsage(nodes map[string]NodeStatus, memThreshold int64) []NodeStatus {
	res := []NodeStatus{}
	for _, node := range nodes {
		if node.memReq > memThreshold {
			res = append(res, node)
		}
	}
	return res
}

func filterPodsByNodeName(nodeName string, resources []*ContainerResources) []*ContainerResources {
	res := []*ContainerResources{}
	for _, r := range resources {
		if r.NodeName == nodeName {
			res = append(res, r)
		}
	}
	return res
}

func GetPodsByUsage(nodeName string, resources []*ContainerResources) []*ContainerResources {
	field := "MemReq"
	reverse := false
	filteredResources := filterPodsByNodeName(nodeName, resources)
	sort.Slice(filteredResources, func(i, j int) bool {
		return cmp(filteredResources, field, i, j, reverse)
	})
	return filteredResources
}

// }
// get most used node
// get resources from most used node
// sort container resources by mem usage
// access: resource.MemReq; resource.PercentMemoryReq,PercentMemoryLimit
