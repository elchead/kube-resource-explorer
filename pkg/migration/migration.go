package migration

import (
	"fmt"
	"io/ioutil"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
)

type Migration struct {
	Pod  string
	Node string
}

func New(pod, node string) *Migration {
	return &Migration{pod, node}
}

func (m Migration) GetManifest() runtime.Object {
	// read yaml
	b, err := ioutil.ReadFile("pod_checkpoint.yaml")
	if err != nil {
		log.Fatal(err)
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, groupVersionKind, err := decode(b, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(groupVersionKind)
	fmt.Println(obj)
	return obj.(*v1.Pod) // do clonePod jkub?
	// decode
	// update values
	//
}
