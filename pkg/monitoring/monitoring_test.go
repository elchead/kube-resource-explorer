package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestGetPodMemory(t *testing.T) {
	// podName := ""
	token := "L2gdmNenL3F9KTGskpNtnPY4wfhVfknn" //"1K8JNp8pG7DJmRc4fhIMpy_A7eiNmesLJqx5U7JY1OJSwswel3DM5Ym-H_oUZ8bypSm609HqNhPgwn4l35OYRw==" //"83p4dTYEpqqctlgzFXMIWD9VtCpfKp2mZwwQoCkdZVTD7ItTL8biFwpmdF-JVZnU5u2YK3S5Yg0C6rfzbrp6ZA=="
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
	fmt.Println("Hi")
	assert.Equal(t, "", result)
}
