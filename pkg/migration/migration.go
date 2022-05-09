package migration

import (
	"fmt"
	"os/exec"
)

const kubeconfig = "/home/adrian/config"

type Migration struct {
	Pod       string
	Namespace string
	ScriptDir string
}

func New(pod, namespace string) *Migration {
	return &Migration{Pod: pod, Namespace: namespace, ScriptDir: "/home/adrian/job-scheduler"}
}

// func (m Migration) GetManifest() runtime.Object {
// 	// read yaml
// 	b, err := ioutil.ReadFile("pod_checkpoint.yaml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	decode := scheme.Codecs.UniversalDeserializer().Decode
// 	obj, groupVersionKind, err := decode(b, nil, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(groupVersionKind)
// 	fmt.Println(obj)
// 	return obj.(*v1.Pod) // do clonePod jkub?
// 	// decode
// 	// update values
// 	//
// }

func (m Migration) Migrate() error {
	cmd := exec.Command("/bin/sh", "./tpod_checkpoint.sh")
	cmd.Env = []string{fmt.Sprintf("WORKER=%s", m.Pod), fmt.Sprintf("NS=%s", m.Namespace), fmt.Sprintf("KUBECONFIG=%s", kubeconfig)}
	cmd.Dir = m.ScriptDir
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return err
}
