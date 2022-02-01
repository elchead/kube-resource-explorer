package kube

import (
	"context"
	"sort"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// mem usage in Gi
type NodeStatus struct {
	name     string
	memReq   int64
	memLimit int64
}

// TODO check sort
func GetNodesByUsage(nodes *v1.NodeList, metricsclient *versioned.Clientset) []NodeStatus {
	res := []NodeStatus{}
	for _, node := range nodes.Items {
		mc, err := metricsclient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.GetName(), metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		usedMem := mc.Usage.Memory().ScaledValue(resource.Scale(9))
		lim := node.Status.Allocatable.Memory().ScaledValue(resource.Scale(9))
		node := NodeStatus{node.GetName(), usedMem, lim}
		res = append(res, node)
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
