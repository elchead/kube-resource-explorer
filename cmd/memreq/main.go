package main

import (
	// "net/http"
	"fmt"
	"log"
	"os"

	"github.com/elchead/kube-resource-explorer/pkg/migration"
	"github.com/elchead/kube-resource-explorer/pkg/monitoring"
	"github.com/joho/godotenv"
	// "github.com/elchead/kube-resource-explorer/pkg/kube"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// "k8s.io/client-go/rest"
)

var token string

func init() {

	err := godotenv.Load("/home/adrian/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token = os.Getenv("INFLUXDB_TOKEN")
}

func main() {

	// var (
	// 	namespace = flag.String("namespace", "playground", "filter by namespace (defaults to all)")
	// 	isLocal   = flag.Bool("isLocal", true, "otherwise use in cluster config")
	// 	config    *rest.Config
	// )
	// flag.Parse()
	// if *isLocal {
	// 	config = kube.GetConfig()

	// } else {
	// 	var err error
	// 	config, err = rest.InClusterConfig()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println(namespace, config)
	// k := kube.GetClient(config)

	url := "https://westeurope-1.azure.cloud2.influxdata.com"
	org := "stobbe.adrian@gmail.com"
	sut := monitoring.New(url, token, org, "default")
	node := "zone2"
	namespace := "playground"
	// pod := "counterten"
	usage, err := sut.GetFreeMemoryNode(node)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Free memory of", node, usage, "percent")
	podMems, err := sut.GetPodMemories(node)
	if err != nil {
		log.Fatal(err)
	}
	for pod := range podMems {
		fmt.Println("Pod", pod, "uses", podMems[pod], "GB")
		fmt.Println("Migrating", pod)
		sut := migration.New(pod, namespace)
		err := sut.Migrate()
		if err != nil {
			log.Fatal(err)
		}
	}

	// result, err := sut.GetPodMemory("counterten", "counterten", "-1ms")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for result.Next() {
	// 	// Notice when group key has changed
	// 	if result.TableChanged() {
	// 		fmt.Printf("table: %s\n", result.TableMetadata().String())
	// 	}
	// 	// Access data
	// 	fmt.Printf("value: %v\n", result.Record().Value())
	// }
	// ticker := time.NewTicker(3 * time.Second)
	// quit := make(chan struct{})
	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		nodes := kube.GetNodesByUsage(kube.GetNodesListAndMetrics(config))
	// 		node, err := kube.GetWorkerNode(nodes)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		memAlloc := nodes[node].MemAlloc

	// 		pods, _ := kube.GetPodsByNode(*k.Clientset, node, *namespace)
	// 		memReqs, memLim := kube.GetPodsTotalMemRequestsAndLimits(pods.Items)
	// 		fractionMemoryReq := float64(memReqs) / float64(memAlloc) * 100

	// 		memReqThresholdPercent := 0.1
	// 		if fractionMemoryReq > memReqThresholdPercent {
	// 			fmt.Printf("Memory request above %f %%!\n", memReqThresholdPercent)
	// 			fmt.Println("Memreq (Gi)", memReqs, "\nMemlim (Gi)", memLim, "\nMemAlloc (Gi)", memAlloc, "\nFrac(%)", fractionMemoryReq)

	// 			podName := pods.Items[0].Name
	// 			fmt.Printf("Checkpoint pod %s\n", podName)
	// 			RequestCheckpointing(podName)

	// 		}
	// 	case <-quit:
	// 		ticker.Stop()
	// 		return
	// 	}
	// }
}

// func RequestCheckpointing(podName string) *http.Response {
// 	url := fmt.Sprintf("http://%s.subdomain:%d/checkpoint", podName, 5747)
// 	// url := fmt.Sprintf("http://subdomain:%d/checkpoint", 5747)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Request Status: ", resp.StatusCode)
// 	return resp
// }
