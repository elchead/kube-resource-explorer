package main

import (
	"strings"
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/kube"
	"github.com/stretchr/testify/assert"
)

func TestGetNodeName(t *testing.T) {
	config := kube.GetConfig()
	nodes := kube.GetNodesByUsage(kube.GetNodesListAndMetrics(config))
	name, err := kube.GetWorkerNode(nodes)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(name, "opt"))
}

func TestGetPodsOfNode(t *testing.T) {
	config := kube.GetConfig()
	k := kube.GetClient(config)
	namespace := "playground"
	node := "shoot--oaas-dev--playground-worker-opt-z2-6858f-bsh4v"
	pods, _ := kube.GetPodsByNode(*k.Clientset, node, namespace)
	assert.NotEmpty(t, pods.Items)
}
