package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/morganhoward/kube-resource-explorer/pkg/kube"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
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
		namespace  = flag.String("namespace", "", "filter by namespace (defaults to all)")
		sort       = flag.String("sort", "CpuReq", "field to sort by")
		reverse    = flag.Bool("reverse", false, "reverse sort output")
		csv        = flag.Bool("csv", false, "Export results to csv file")
		version    = flag.Bool("version", false, "show binary version")
		kubeconfig *string
	)

	if kubeenv := os.Getenv("KUBECONFIG"); kubeenv != "" {
		kubeconfig = flag.String("kubeconfig", kubeenv, "absolute path to the kubeconfig file")
	} else if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	if *version {
		fmt.Printf("Build information {OS:%q Arch:%q GitCommit:%q, GoVersion:%q}\n", runtime.GOOS, runtime.GOARCH, GitCommit, runtime.Version())
		os.Exit(0)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k := kube.NewKubeClient(clientset)

	r := kube.ContainerResources{}

	if !r.Validate(*sort) {
		fmt.Printf("\"%s\" is not a valid field. Possible values are:\n\n%s\n", *sort, strings.Join(kube.GetFields(r), ", "))
		os.Exit(1)
	}

	k.ResourceUsage(*namespace, *sort, *reverse, *csv)
}
