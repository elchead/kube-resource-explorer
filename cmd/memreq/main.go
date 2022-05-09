package main

import (
	// "net/http"

	"log"
	"os"
	"time"

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
	client := monitoring.NewWithTime(url, token, org, "default", "-5h")
	namespace := "playground"
	ctrl := monitoring.NewController(client)

	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			migs, _ := ctrl.GetMigrations()
			migration.Migrate(migs, namespace)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
