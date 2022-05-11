package monitoring

import (
	"github.com/elchead/kube-resource-explorer/pkg/migration"
)

type Controller struct {
	Client  Clienter
	Policy  MigrationPolicy
	Cluster Cluster
}

func NewController(client Clienter) *Controller {
	return NewControllerWithPolicy(client, ThresholdPolicy{20.})
}

func NewControllerWithPolicy(client Clienter, policy MigrationPolicy) *Controller {
	return &Controller{client, policy, Cluster{NbrNodes: 2, NodeGb: 450., TargetFreePercent: 20.}}
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

type NodeFullError struct{}

func (m *NodeFullError) Error() string {
	return "nodes are full. no place to migrate"
}

func (c Controller) GetMigrations() (migrations []migration.MigrationCmd, err error) {
	criticalNode := c.Policy.GetCriticalNodes(c.Client)
	nodeMemApp, _ := c.Client.GetFreeMemoryOfNodes()
	if !c.Cluster.isSpaceAvailable(criticalNode) {
		return migrations, &NodeFullError{}
	}
	for _, node := range criticalNode {
		podMems, err := c.Client.GetPodMemories(node)
		if err != nil {
			return migrations, err
		}
		pod := GetMaxPod(podMems)
		podMem := podMems[pod]
		if c.Cluster.enoughSpaceAvailableOn(node, podMem, nodeMemApp) != "" {
			migrations = append(migrations, migration.MigrationCmd{Pod: pod, Usage: podMem})
		} else {
			return migrations, &NodeFullError{}
		}
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
