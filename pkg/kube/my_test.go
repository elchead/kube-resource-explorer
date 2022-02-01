package kube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func TestNodesByMem(t *testing.T) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	metricsclient, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k := NewKubeClient(clientset)

	// resources, err := k.ContainerResources(namespace)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// capacity, err := k.ClusterCapacity()
	// if err != nil {
	// 	panic(err.Error())
	// }

	nodes, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	res := GetNodesByUsage(nodes, metricsclient)
	fmt.Println(res)
	assert.NotEqual(t, len(res), 0)
}

func TestGetPodsByUsage(t *testing.T) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k := NewKubeClient(clientset)

	namespace := ""
	resources, err := k.ContainerResources(namespace)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeName := nodes.Items[0].Name
	res := GetPodsByUsage(nodeName, resources)
	fmt.Println(*res[0])
	assert.False(t, true)
}
