package monitoring

type Cluster struct {
	NbrNodes          int
	NodeGb            float64
	TargetFreePercent float64
}

func (c Cluster) isSpaceAvailable(nodes []string) bool {
	return len(nodes) != c.NbrNodes
}

func (c Cluster) getFreePercent(freeNodeGb float64) float64 {
	return freeNodeGb / c.NodeGb * 100.
}

func (c Cluster) enoughSpaceAvailableOn(originalNode string, podMemory float64, nodes NodeFreeMemMap) string {
	for node, free_percent := range nodes {
		if node != originalNode {
			freeGb := free_percent / 100. * c.NodeGb
			newFreeGb := freeGb - podMemory
			if c.getFreePercent(newFreeGb) > c.TargetFreePercent {
				return node
			}
		}
	}
	return ""
}
