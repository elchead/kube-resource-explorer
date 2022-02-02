package kube

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPodsByNode(t *testing.T) {
	node := "shoot--oaas-live--play-worker-opt-z3-5b545-g7htk"
	config := GetConfig()
	k := GetClient(config)
	res, _ := GetPodsByNode(*k.Clientset, node, "play")
	fmt.Println(res)
	assert.Equal(t, true, false)

}

func TestGetNodeMemReq(t *testing.T) {
	node := "shoot--oaas-live--play-worker-opt-z3-5b545-g7htk"
	config := GetConfig()
	k := GetClient(config)
	res, _ := GetPodsByNode(*k.Clientset, node, "play")
	memReqs, memLim := GetPodsTotalMemRequestsAndLimits(res.Items)
	assert.Equal(t, 0, memReqs)
	assert.Equal(t, 0, memLim)

}
func TestGetNodeWithXUsage(t *testing.T) {
	// config := getConfig()
	// res := GetNodesByUsage(getNodesListAndMetrics(config))
	nodes := []NodeStatus{{"low", 1, 10}, {"high", 6, 10}}
	res := FilterNodesByUsage(nodes, 5)
	assert.Contains(t, res, NodeStatus{"high", 6, 10})
	assert.NotContains(t, res, NodeStatus{"low", 1, 10})
}

func TestFindIntensivePodOnCriticalNode(t *testing.T) {
	config := GetConfig()
	nodes := GetNodesByUsage(getNodesListAndMetrics(config))
	fmt.Println(nodes)
	res := FilterNodesByUsage(nodes, 5)
	criticalNode := res[0]

	k := GetClient(config)
	namespace := ""
	resources, err := k.ContainerResources(namespace)
	if err != nil {
		panic(err.Error())
	}
	pods := GetPodsByUsage(criticalNode.name, resources)
	assert.Equal(t, "", pods[0].Name)
}

func TestNodesByMem(t *testing.T) {
	config := GetConfig()
	res := GetNodesByUsage(getNodesListAndMetrics(config))
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
