package monitoring

import (
	"context"
	"errors"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
	client   influxdb2.Client
	queryAPI api.QueryAPI
	bucket   string
}

func NewLocal(token, org, bucket string) *Client {
	return New("http://localhost:8081", token, org, bucket)
}

func New(serviceUrl, token, org, bucket string) *Client {
	client := influxdb2.NewClientWithOptions(serviceUrl, token, influxdb2.DefaultOptions())
	return &Client{client, client.QueryAPI(org), bucket}
}

func (c *Client) Query(query string) (*api.QueryTableResult, error) {
	return c.queryAPI.Query(context.Background(), query)
}

func (c *Client) GetPodMemory(podName, containerName, time string) (*api.QueryTableResult, error) {
	query := fmt.Sprintf(`from(bucket: "%s") 
	|> range(start: %s)
	|> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
	|> filter(fn: (r) => r["_field"] == "memory_usage_bytes")
	|> filter(fn: (r) => r["pod_name"] == "%s")
	|> filter(fn: (r) => r["container_name"] == "%s" )
  |> yield(name: "mean")`, c.bucket, time, podName, containerName)
	return c.Query(query) // default container: worker
}

func (c *Client) GetPodMemorySlope(podName, time, slopeWindow string) (float64, error) {
	query := fmt.Sprintf(`import "experimental/aggregate" from(bucket: "%s") 
  |> range(start: %s)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
  |> filter(fn: (r) => r["_field"] == "memory_usage_bytes")
  |> filter(fn: (r) => r["pod_name"] == "%s")
  |> filter(fn: (r) => r["container_name"] == "worker")
  |> aggregate.rate(every: %s, unit: 1m, groupColumns: ["tag1", "tag2"])
  |> mean()`, c.bucket, time, podName, slopeWindow)
	res, err := c.Query(query)
	if res.Next() && err == nil {
		num := res.Record().Value()
		if val, ok := num.(float64); ok {
			return val, nil
		} else {
			return -1., errors.New("conversion error")
		}
	}
	return -1., err
}

func (c *Client) GetFreeMemoryNode(nodeName string) (float64, error) {
	time := "-1m"
	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: %s)
	|> filter(fn: (r) => r["_measurement"] == "mem")
	|> filter(fn: (r) => r["_field"] == "available_percent")
	|> filter(fn: (r) => r["host"] == "%s")
	|> last()`, c.bucket, time, nodeName)
	res, err := c.Query(query)
	if err == nil && res.Next() {
		num := res.Record().Value()
		// fmt.Println(num)
		// fmt.Println(res.Record())
		if val, ok := num.(float64); ok {
			return val, nil
		}
	}
	return -1., err
}

type NodeMemMap map[string]float64

func (c *Client) GetFreeMemoryOfNodes() (NodeMemMap, error) {
	time := "-1m"
	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: %s)
	|> filter(fn: (r) => r["_measurement"] == "mem")
	|> filter(fn: (r) => r["_field"] == "available_percent")
	|> last()`, c.bucket, time)
	res, err := c.Query(query)

	mp := make(NodeMemMap)
	for err == nil && res.Next() {
		table := res.Record()
		node := table.ValueByKey("host").(string)
		available_percent := table.Value().(float64)
		mp[node] = available_percent
	}
	return mp, err
}
