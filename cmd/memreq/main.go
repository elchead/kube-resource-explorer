package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/elchead/kube-resource-explorer/pkg/kube"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
)

var GitCommit string

const node = "shoot--oaas-dev--playground-worker-opt-z2-6858f-bsh4v"

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {

	var (
		namespace = flag.String("namespace", "playground", "filter by namespace (defaults to all)")
		isLocal   = flag.Bool("isLocal", true, "otherwise use in cluster config")
		config    *rest.Config
	)
	// 	kubeconfig *string
	// )

	// if kubeenv := os.Getenv("KUBECONFIG"); kubeenv != "" {
	// 	kubeconfig = flag.String("kubeconfig", kubeenv, "absolute path to the kubeconfig file")
	// } else if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }

	flag.Parse()
	if *isLocal {
		config = kube.GetConfig()

	} else {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}
	k := kube.GetClient(config)
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			// do stuff
			// nodes := kube.GetNodesByUsage(kube.GetNodesListAndMetrics(config))
			// memAlloc := nodes[node].MemAlloc
			printUsage(k, node, namespace, config)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func printUsage(k *kube.KubeClient, node string, namespace *string, config *rest.Config) {
	res, err := kube.GetPodsByNode(*k.Clientset, node, *namespace)
	if err != nil {
		panic(err)
	}
	if len(res.Items) == 0 {
		panic(fmt.Errorf("No pods found for %s, ns:%s", node, *namespace))
	}
	memReqs, memLim := kube.GetPodsTotalMemRequestsAndLimits(res.Items)
	// if memR
	nodes := kube.GetNodesByUsage(kube.GetNodesListAndMetrics(config))
	memAlloc := nodes[node].MemAlloc
	fractionMemoryReq := float64(memReqs) / float64(memAlloc) * 100
	if fractionMemoryReq > 0.1 {
		fmt.Println("Above 0.1%!")
		pods, _ := kube.GetPodsByNode(*k.Clientset, node, *namespace)
		podName := pods.Items[0].Name
		fmt.Println("Pod", podName)
		fmt.Println("Send checkpointing command")
		TEST := "test"
		resp := RequestCheckpointing(TEST)
		fmt.Println("Status: ", resp.StatusCode)

	}
	fmt.Println("Memreq", memReqs, "\nMemlim", memLim, "\nMemAlloc", memAlloc, "\nFrac", fractionMemoryReq)
	fmt.Println(nodes)
}

func RequestCheckpointing(podName string) *http.Response {
	url := fmt.Sprintf("http://%s.subdomain:%d/checkpoint", podName, 5747)
	// url := fmt.Sprintf("http://subdomain:%d/checkpoint", 5747)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp
}
