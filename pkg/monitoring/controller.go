package monitoring

import (
	"github.com/elchead/kube-resource-explorer/pkg/migration"
)

type Controller struct {
	Client  Clienter
	Policy  MigrationPolicy
	Cluster Cluster
}

type Cluster struct {
	NbrNodes          int
	NodeSize          float64
	TargetFreePercent float64
}

func (c Cluster) isSpaceAvailable(nodes []string) bool {
	if len(nodes) == c.NbrNodes {
		return false
	}
	return true
}

func (c Cluster) getFreePercent(freeNodeGb float64) float64 {
	return freeNodeGb / c.NodeSize * 100.
}

func (c Cluster) enoughSpaceAvailableOn(originalNode string, podMemory float64, nodes NodeFreeMemMap) string {
	for node, free_percent := range nodes {
		if node != originalNode {
			freeGb := free_percent / 100. * c.NodeSize
			newFreeGb := freeGb - podMemory
			if c.getFreePercent(newFreeGb) > c.TargetFreePercent {
				return node
			}
		}
	}
	return ""
}
func NewController(client Clienter) *Controller {
	return NewControllerWithPolicy(client, ThresholdPolicy{20.})
}

func NewControllerWithPolicy(client Clienter, policy MigrationPolicy) *Controller {
	return &Controller{client, policy, Cluster{NbrNodes: 2, NodeSize: 450., TargetFreePercent: 20.}}
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
	nodes := c.Policy.GetCriticalNodes(c.Client)
	// nodeMemApp, _ := c.Client.GetFreeMemoryOfNodes()
	if !c.Cluster.isSpaceAvailable(nodes) {
		return []migration.MigrationCmd{}, &NodeFullError{}
	}
	for _, node := range nodes {
		podMems, err := c.Client.GetPodMemories(node)
		if err != nil {
			return migrations, err
		}
		pod := GetMaxPod(podMems)
		podMem := podMems[pod]
		migrations = append(migrations, migration.MigrationCmd{Pod: pod, Usage: podMem})
		// if c.Cluster.enoughSpaceAvailableOn(node, podMem, nodeMemApp) != "" {
		// }
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
