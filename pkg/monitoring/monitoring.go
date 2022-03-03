package monitoring

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
	client   influxdb2.Client
	queryAPI api.QueryAPI
	bucket   string
}

func New(token, org, bucket string) *Client {
	client := influxdb2.NewClientWithOptions("http://localhost:8081", token, influxdb2.DefaultOptions())
	return &Client{client, client.QueryAPI(org), bucket}
}

func (c *Client) Query(query string) (*api.QueryTableResult, error) {
	return c.queryAPI.Query(context.Background(), query)
}

func (c *Client) GetPodMemory(podName, time string) (*api.QueryTableResult, error) {
	query := fmt.Sprintf(`from(bucket: "%s") 
  |> range(start: %s)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
  |> filter(fn: (r) => r["_field"] == "memory_usage_bytes")
  |> filter(fn: (r) => r["pod_name"] == "%s")
  |> filter(fn: (r) => r["container_name"] == "worker")
  |> yield(name: "mean")`, c.bucket, time, podName)
	return c.Query(query)
}
