package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/monitoring"
	"github.com/stretchr/testify/assert"
)

const token = "L2gdmNenL3F9KTGskpNtnPY4wfhVfknn"

func TestGetPodMemory(t *testing.T) {
	sut := monitoring.NewLocal(token, "influxdata", "default")
	result, err := sut.GetPodMemory("acounterten", "counterten", "-1m")
	assert.NoError(t, err)
	for result.Next() {
		// Notice when group key has changed
		// if result.TableChanged() {
		// 	fmt.Printf("table: %s\n", result.TableMetadata().String())
		// }
		fmt.Println(result.Record())
		// Access data
		// fmt.Printf("value: %v\n", result.Record().Value())
	}

	assert.True(t, false)
}
func TestPodMemoryOfNode(t *testing.T) {
	sut := monitoring.NewLocal(token, "influxdata", "default")
	res, err := sut.GetPodMemories("shoot--oaas-dev--playground-worker-opt-z2-6bf98-9dv44")
	fmt.Println(res)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.NotEqual(t, -1., res["o10n-worker-s-ll8rn-75t9n"])
}

func TestFreeMemoryNode(t *testing.T) {
	sut := monitoring.NewLocal(token, "influxdata", "default")
	res, err := sut.GetFreeMemoryNode("shoot--oaas-dev--playground-worker-opt-z2-6bf98-9dv44")
	assert.NoError(t, err)
	assert.Equal(t, -1., res)
}

func TestFreeMemoryOfNodes(t *testing.T) {
	sut := monitoring.NewLocal(token, "influxdata", "default")
	res, err := sut.GetFreeMemoryOfNodes()
	fmt.Println(res)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, -1., res["shoot--oaas-dev--playground-worker-opt-z2-6bf98-9dv44"])
}

func TestGetPodMemorySlope(t *testing.T) {
	sut := monitoring.NewLocal(token, "influxdata", "default")
	result, err := sut.GetPodMemorySlope("o10n-worker-s-qxw8k-pbclk", "-3m", "1m")
	assert.NoError(t, err)
	assert.Equal(t, 0., result)
	// for result.Next() {
	// 	// Notice when group key has changed
	// 	if result.TableChanged() {
	// 		fmt.Printf("table: %s\n", result.TableMetadata().String())
	// 	}
	// 	// Access data
	// 	fmt.Printf("value: %v\n", result.Record().Value())
	// }
	// assert.Equal(t, "", result.Record().Value())
}
