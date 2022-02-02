package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/elchead/kube-resource-explorer/pkg/kube"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
)

var GitCommit string

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {

	var (
		namespace = flag.String("namespace", "play", "filter by namespace (defaults to all)")
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
	node := "shoot--oaas-live--play-worker-opt-z3-5b545-g7htk"
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
	res, _ := kube.GetPodsByNode(*k.Clientset, node, *namespace)
	memReqs, memLim := kube.GetPodsTotalMemRequestsAndLimits(res.Items)
	fmt.Println("Memreq", memReqs, "\nMemlim", memLim)
	nodes := kube.GetNodesByUsage(kube.GetNodesListAndMetrics(config))
	fmt.Println(nodes)
}
