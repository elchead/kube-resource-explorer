package monitoring_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/monitoring"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var token string
var sut *monitoring.Client

const url = "https://westeurope-1.azure.cloud2.influxdata.com"
const org = "stobbe.adrian@gmail.com"
const node = "zone2"
const pod = "counterten"

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token = os.Getenv("INFLUXDB_TOKEN")
	sut = monitoring.New(url, token, org, "default")
}

func TestGbConversion(t *testing.T) {
	assert.Equal(t, 192.7, monitoring.ConvertToGb(206912778240))
}
func TestPodMemoryOfNode(t *testing.T) {
	res, err := sut.GetPodMemoriesFromContainer(node, pod)
	fmt.Println(res)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.NotEqual(t, -1., res[pod])
	assert.NotEqual(t, 0., res[pod])
}

func TestFreeMemoryNode(t *testing.T) {
	res, err := sut.GetFreeMemoryNode(node)
	assert.NoError(t, err)
	assert.NotEqual(t, -1., res)
}

func TestFreeMemoryOfNodes(t *testing.T) {
	res, err := sut.GetFreeMemoryOfNodes()
	fmt.Println(res)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.NotEqual(t, -1., res[node])
}

func TestGetPodMemorySlope(t *testing.T) {
	result, err := sut.GetPodMemorySlopeFromContainer(pod, "counterten", "-3m", "1m")
	assert.NoError(t, err)
	assert.Equal(t, 0., result)
}
