package monitoring

type Controller struct {
	Client Clienter
	Policy MigrationPolicy
}

func NewController(client Clienter) *Controller {
	return &Controller{client, ThresholdPolicy{}}
}

type MigrationCmd struct {
	Pod string
}

type ThresholdPolicy struct{}

func (t ThresholdPolicy) GetCriticalNodes(c Clienter) (criticalNodes []string) {
	nodes, _ := c.GetFreeMemoryOfNodes()
	for node := range nodes {
		free := nodes[node]
		if free < 20. {
			criticalNodes = append(criticalNodes, node)
		}
	}
	return criticalNodes
}

type MigrationPolicy interface {
	GetCriticalNodes(Clienter) []string
}

// TODO remove
func (c Controller) GetCriticalNodes() []string {
	return c.Policy.GetCriticalNodes(c.Client)
}

func (c Controller) GetMigrations() (migrations []MigrationCmd, err error) {
	nodes := c.Policy.GetCriticalNodes(c.Client)
	for _, node := range nodes {
		podMems, err := c.Client.GetPodMemories(node)
		if err != nil {
			return migrations, err
		}
		pod := GetMaxPod(podMems)
		migrations = append(migrations, MigrationCmd{pod})
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
