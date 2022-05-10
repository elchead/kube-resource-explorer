package monitoring

import (
	"github.com/elchead/kube-resource-explorer/pkg/migration"
)

type Controller struct {
	Client Clienter
	Policy MigrationPolicy
}

func NewController(client Clienter) *Controller {
	return &Controller{client, ThresholdPolicy{20.}}
}

type ThresholdPolicy struct {
	ThresholdFreePercent float64
}

func (t ThresholdPolicy) GetCriticalNodes(c Clienter) (criticalNodes []string) {
	nodes, _ := c.GetFreeMemoryOfNodes()
	for node, free := range nodes {
		if free < t.ThresholdFreePercent {
			criticalNodes = append(criticalNodes, node)
		}
	}
	return criticalNodes
}

type MigrationPolicy interface {
	GetCriticalNodes(Clienter) []string
}

func (c Controller) GetMigrations() (migrations []migration.MigrationCmd, err error) {
	nodes := c.Policy.GetCriticalNodes(c.Client)
	for _, node := range nodes {
		podMems, err := c.Client.GetPodMemories(node)
		if err != nil {
			return migrations, err
		}
		pod := GetMaxPod(podMems)
		migrations = append(migrations, migration.MigrationCmd{Pod: pod, Usage: podMems[pod]})
	}
	return migrations, nil
}

func GetMaxPod(pods PodMemMap) (pod string) {
	max := -1.
	for p, mem := range pods {
		if mem > max {
			max = mem
			pod = p
		}
	}
	return pod
}
