package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestGetPodMemory(t *testing.T) {
	token := "L2gdmNenL3F9KTGskpNtnPY4wfhVfknn"
	sut := monitoring.New(token, "influxdata", "default")
	result, err := sut.GetPodMemory("worker-l-vj9vv-8p5wr", "-1h")
	assert.NoError(t, err)
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
	}
	assert.Equal(t, "", result.Record().Value())
}

func TestGetPodMemorySlope(t *testing.T) {
	token := "L2gdmNenL3F9KTGskpNtnPY4wfhVfknn"
	sut := monitoring.New(token, "influxdata", "default")
	result, err := sut.GetPodMemorySlope("worker-m-jcbxp-hk75j", "-1h", "20m")
	assert.NoError(t, err)
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
	}
	assert.Equal(t, "", result.Record().Value())
}
