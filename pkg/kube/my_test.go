package kube

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const node = "shoot--oaas-dev--playground-worker-opt-z2-6858f-bsh4v"
const namespace = "playground"

func TestGetPodsByNode(t *testing.T) {
	config := GetConfig()
	k := GetClient(config)
	res, _ := GetPodsByNode(*k.Clientset, node, namespace)
	fmt.Println(res)
	assert.Equal(t, true, false)

}

func TestGetNodeMemReq(t *testing.T) {
	config := GetConfig()
	k := GetClient(config)
	res, _ := GetPodsByNode(*k.Clientset, node, namespace)
	memReqs, memLim := GetPodsTotalMemRequestsAndLimits(res.Items)
	assert.Equal(t, 0, memReqs)
	assert.Equal(t, 0, memLim)

}
func TestGetNodeWithXUsage(t *testing.T) {
	// config := getConfig()
	// res := GetNodesByUsage(getNodesListAndMetrics(config))
	nodes := map[string]NodeStatus{"low": {"low", 1, 10}, "high": {"high", 6, 10}}
	res := FilterNodesByUsage(nodes, 5)
	assert.Contains(t, res, NodeStatus{"high", 6, 10})
	assert.NotContains(t, res, NodeStatus{"low", 1, 10})
}

func TestFindIntensivePodOnCriticalNode(t *testing.T) {
	config := GetConfig()
	nodes := GetNodesByUsage(GetNodesListAndMetrics(config))
	fmt.Println(nodes)
	res := FilterNodesByUsage(nodes, 1)
	criticalNode := res[0]

	k := GetClient(config)
	resources, err := k.ContainerResources(namespace)
	if err != nil {
		panic(err.Error())
	}
	pods := GetPodsByUsage(criticalNode.name, resources)
	assert.Equal(t, "", pods[0].Name)
}

func TestNodesByMem(t *testing.T) {
	config := GetConfig()
	res := GetNodesByUsage(GetNodesListAndMetrics(config))
	fmt.Println(res)
	assert.NotEqual(t, len(res), 0)
}

func TestGetPodsByUsage(t *testing.T) {
	config := GetConfig()
	k := GetClient(config)

	namespace := ""
	resources, err := k.ContainerResources(namespace)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := k.NodeList()
	if err != nil {
		panic(err.Error())
	}
	nodeName := nodes.Items[0].Name
	res := GetPodsByUsage(nodeName, resources)
	fmt.Println(*res[0])
	assert.NotEmpty(t, res)
}
