# My

`memreq` observes the node usage and triggers checkpointing.

## Cluster Deployment
1. For API permissions, add service-account in worker namespace to `cluster-admin` Clusterrolebinding:

    ```
    k edit clusterrolebinding cluster-admin
    // add
    20 subjects:
    21 - apiGroup: rbac.authorization.k8s.io
    22   kind: Group
    23   name: system:masters
    24 - kind: ServiceAccount
    25   name: default
    26   namespace: playground
    ```
2. Deploy pod service `kaf pod-svc.yaml`
3. Deploy worker job (modify node name for affinity!!): `kaf worker.yaml`
4. Wait about 1min for worker job to deploy (but shuts down after a few minutes!!)
5. [Check availability with : `wget http://test.subdomain:5747/checkpoint`]
6. Deploy `kaf memreq.yaml`

## Development

`make docker`: build docker image and push to repo
`make my`: build `memreq`

## Checkpointing
Provide pod name:
`WORKER=worker-l-x2bzh-z29tl ./checkpoint.sh`
---
# Resource Explorer

Note: This fork doesn't use Google cloud resources and has removed this functionality.

[![CircleCI](https://circleci.com/gh/dabeck/kube-resource-explorer/tree/master.svg?style=svg)](https://circleci.com/gh/dabeck/kube-resource-explorer/tree/master)

Explore your kube resource usage and allocation.

* Display historical statistical resource usage from StackDriver.

  <https://github.com/kubernetes/kubernetes/issues/55046>

* List resource QoS allocation to pods in a cluster. Inspired by:

  <https://github.com/kubernetes/kubernetes/issues/17512>

## Usage

### Command Line Options

* `-namespace` - Limit the query to the specified namespace (defaults to all)
* `-sort` - Field to sort by
* `-reverse` - Reserve the sort order
* `-csv` - Export results to CSV file
* `-version` - Print the binary version

### Run

```sh
make run
```

### Build

```sh
make build
```

### Build + Install

```sh
make
```

## Example output

Show aggregate resource requests and limits. This is the same information
displayed by `kubectl describe nodes` but in a easier to view format.

```sh
$ kube-resource-explorer -reverse -sort MemReq
Node        Namespace                 Name                                                                       CpuReq       CpuReq%  CpuLimit     CpuLimit%  MemReq          MemReq%  MemLimit        MemLimit%  Pod Age
----        ---------                 ----                                                                       ------       -------  --------     ---------  ------          -------  --------        ---------  -------
local-node  cattle-monitoring-system  pushprox-k3s-server-proxy-f4f5d4874-689xb/pushprox-proxy                   0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m48s
local-node  cattle-fleet-system       fleet-controller-974d9cc9f-csf66/fleet-controller                          0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h41m49s
local-node  cattle-fleet-system       gitjob-5778966b7c-jmtkr/gitjob                                             0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h41m49s
local-node  cattle-logging-system     rancher-logging-d9bf878c6-7nqxg/rancher-logging                            0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h28m55s
local-node  istio-system              kiali-7c4c559b9f-7vzdg/kiali                                               0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h28m13s
local-node  default                   gentoo-7b759fdbf8-7gcgh/container-0                                        0m           0%       0m           0%         0Mi             0%       0Mi             0%         1318h55m51s
local-node  cattle-logging-system     rancher-logging-fluentd-0/config-reloader                                  0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h28m33s
local-node  cattle-logging-system     rancher-logging-k3s-journald-aggregator-2rbq2/fluentbit                    0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h28m55s
local-node  cattle-system             rancher-webhook-7f84b74ddb-rrj7w/rancher-webhook                           0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h41m29s
local-node  cattle-monitoring-system  rancher-monitoring-prometheus-adapter-77568b975-2fxq6/prometheus-adapter   0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m48s
local-node  cattle-monitoring-system  rancher-monitoring-grafana-8686947796-zw2xk/grafana-proxy                  0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m48s
local-node  cattle-monitoring-system  rancher-monitoring-grafana-8686947796-zw2xk/grafana-sc-dashboard           0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m48s
local-node  cattle-monitoring-system  prometheus-rancher-monitoring-prometheus-0/prometheus-proxy                0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m40s
local-node  cattle-monitoring-system  pushprox-k3s-server-client-cfqdm/pushprox-client                           0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h30m48s
local-node  cattle-fleet-system       fleet-agent-b94869475-bfb9v/fleet-agent                                    0m           0%       0m           0%         0Mi             0%       0Mi             0%         3402h41m20s
local-node  cattle-monitoring-system  rancher-monitoring-prometheus-node-exporter-dd968/node-exporter            100m         1%       200m         3%         30Mi            0%       50Mi            0%         3402h30m48s
local-node  cattle-logging-system     rancher-logging-fluentbit-42ks2/fluent-bit                                 100m         1%       200m         3%         47Mi            0%       95Mi            0%         3402h28m33s
local-node  cattle-monitoring-system  alertmanager-rancher-monitoring-alertmanager-0/config-reloader             100m         1%       100m         1%         50Mi            0%       50Mi            0%         3402h30m40s
local-node  cattle-monitoring-system  prometheus-rancher-monitoring-prometheus-0/config-reloader                 100m         1%       100m         1%         50Mi            0%       50Mi            0%         3402h30m40s
local-node  kube-system               coredns-7448499f4d-rq7mp/coredns                                           100m         1%       0m           0%         70Mi            0%       170Mi           0%         3402h42m23s
local-node  cattle-logging-system     rancher-logging-fluentd-0/fluentd                                          500m         8%       1000m        16%        95Mi            0%       381Mi           1%         3402h28m33s
local-node  cattle-monitoring-system  rancher-monitoring-operator-754bcd8cb4-hqpjb/rancher-monitoring            100m         1%       200m         3%         100Mi           0%       500Mi           1%         3402h30m48s
local-node  cattle-monitoring-system  rancher-monitoring-grafana-8686947796-zw2xk/grafana                        100m         1%       200m         3%         100Mi           0%       200Mi           0%         3402h30m48s
local-node  cattle-monitoring-system  alertmanager-rancher-monitoring-alertmanager-0/alertmanager                100m         1%       1000m        16%        100Mi           0%       500Mi           1%         3402h30m40s
local-node  istio-system              istio-egressgateway-9c86c49bb-nth6x/istio-proxy                            100m         1%       2000m        33%        128Mi           0%       1024Mi          3%         3402h27m24s
local-node  istio-system              istio-ingressgateway-5d84c54d96-sk4vt/istio-proxy                          100m         1%       2000m        33%        128Mi           0%       1024Mi          3%         3402h27m24s
local-node  cattle-monitoring-system  rancher-monitoring-kube-state-metrics-744b9448f4-gbw5j/kube-state-metrics  100m         1%       100m         1%         130Mi           0%       200Mi           0%         3402h30m48s
local-node  cattle-monitoring-system  prometheus-rancher-monitoring-prometheus-0/prometheus                      750m         12%      1000m        16%        1750Mi          5%       2500Mi          7%         3402h30m40s
local-node  istio-system              istiod-5f95cf9cbf-248nm/discovery                                          500m         8%       0m           0%         2048Mi          6%       0Mi             0%         3402h27m36s
----        ---------                 ----                                                                       ------       -------  --------     ---------  ------          -------  --------        ---------  -------
Total                                                                                                            2850m/6000m  47%      8100m/6000m  135%       4827Mi/31666Mi  15%      6744Mi/31666Mi  21%

```
