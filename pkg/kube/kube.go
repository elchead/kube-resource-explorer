package kube

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubeClient struct {
	Clientset *kubernetes.Clientset
}

func NewKubeClient(clientset *kubernetes.Clientset) *KubeClient {
	return &KubeClient{Clientset: clientset}
}

func (c *KubeClient) NodeList() (*corev1.NodeList, error) {
	return c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (k *KubeClient) ActivePods(namespace, nodeName string) ([]corev1.Pod, error) {

	selector := fmt.Sprintf("status.phase!=%s,status.phase!=%s", string(corev1.PodSucceeded), string(corev1.PodFailed))
	if nodeName != "" {
		selector += fmt.Sprintf(",spec.nodeName=%s", nodeName)
	}

	fieldSelector, err := fields.ParseSelector(selector)
	if err != nil {
		return nil, err
	}

	activePods, err := k.Clientset.CoreV1().Pods(
		namespace,
	).List(context.TODO(),
		metav1.ListOptions{FieldSelector: fieldSelector.String()},
	)
	if err != nil {
		return nil, err
	}

	return activePods.Items, err
}

func containerRequestsAndLimits(container *corev1.Container) (reqs corev1.ResourceList, limits corev1.ResourceList) {
	reqs, limits = corev1.ResourceList{}, corev1.ResourceList{}

	for name, quantity := range container.Resources.Requests {
		if _, ok := reqs[name]; ok {
			panic(fmt.Sprintf("Duplicate key: %s", name))
		} else {
			reqs[name] = quantity.DeepCopy()
		}
	}

	for name, quantity := range container.Resources.Limits {
		if _, ok := limits[name]; ok {
			panic(fmt.Sprintf("Duplicate key: %s", name))
		} else {
			limits[name] = quantity.DeepCopy()
		}
	}
	return
}

func NodeCapacity(node *corev1.Node) corev1.ResourceList {
	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}
	return allocatable
}

func (k *KubeClient) NodeResources(namespace, nodeName string) (resources []*ContainerResources, err error) {

	mc := k.Clientset.CoreV1().Nodes()
	node, err := mc.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	activePodsList, err := k.ActivePods(namespace, nodeName)
	if err != nil {
		return nil, err
	}

	capacity := NodeCapacity(node)

	// https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/describe.go#L2970
	for _, pod := range activePodsList {
		// allocatable := node.Status.Capacity
		// if len(node.Status.Allocatable) > 0 {
		// 	allocatable = node.Status.Allocatable
		// }
		req, limit := resourcehelper.PodRequestsAndLimits(&pod)
		for _, container := range pod.Spec.Containers {
			// req, limit := resourcehelper.PodRequestsAndLimits(&pod)
			// cpuReq, cpuLimit, memoryReq, memoryLimit := req[corev1.ResourceCPU], limit[corev1.ResourceCPU], req[corev1.ResourceMemory], limit[corev1.ResourceMemory]
			// req, limit := containerRequestsAndLimits(&container)

			_cpuReq := req[corev1.ResourceCPU]
			cpuReq := NewCpuResource(_cpuReq.MilliValue())

			_cpuLimit := limit[corev1.ResourceCPU]
			cpuLimit := NewCpuResource(_cpuLimit.MilliValue())

			_memoryReq := req[corev1.ResourceMemory]
			memoryReq := NewMemoryResource(_memoryReq.Value())

			_memoryLimit := limit[corev1.ResourceMemory]
			memoryLimit := NewMemoryResource(_memoryLimit.Value())

			// fractionCpuReq := float64(cpuReq.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
			// fractionCpuLimit := float64(cpuLimit.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
			// fractionMemoryReq := float64(memoryReq.Value()) / float64(allocatable.Memory().Value()) * 100
			// fractionMemoryLimit := float64(memoryLimit.Value()) / float64(allocatable.Memory().Value()) * 100

			resources = append(resources, &ContainerResources{
				NodeName:           nodeName,
				Name:               fmt.Sprintf("%s/%s", pod.GetName(), container.Name),
				Namespace:          pod.GetNamespace(),
				CpuReq:             cpuReq,
				CpuLimit:           cpuLimit,
				PercentCpuReq:      cpuReq.calcPercentage(capacity.Cpu()),
				PercentCpuLimit:    cpuLimit.calcPercentage(capacity.Cpu()),
				MemReq:             memoryReq,
				MemLimit:           memoryLimit,
				PercentMemoryReq:   memoryReq.calcPercentage(capacity.Memory()),
				PercentMemoryLimit: memoryLimit.calcPercentage(capacity.Memory()),
				PodAge:             time.Since(pod.GetCreationTimestamp().Time),
			})
		}
	}

	return resources, nil
}

func (k *KubeClient) ContainerResources(namespace string) (resources []*ContainerResources, err error) {
	nodes, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		nodeUsage, err := k.NodeResources(namespace, node.GetName())
		if err != nil {
			return nil, err
		}
		resources = append(resources, nodeUsage...)
	}

	return resources, nil
}

func (k *KubeClient) ClusterCapacity() (capacity corev1.ResourceList, err error) {
	nodes, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	capacity = corev1.ResourceList{}

	for _, node := range nodes.Items {

		allocatable := NodeCapacity(&node)

		for name, quantity := range allocatable {
			if value, ok := capacity[name]; ok {
				value.Add(quantity)
				capacity[name] = value
			} else {
				capacity[name] = quantity.DeepCopy()
			}
		}

	}

	return capacity, nil
}

func (k *KubeClient) ResourceUsage(metricsclient *versioned.Clientset, namespace, sort string, reverse bool, csv bool, nodesOnly bool) {

	resources, err := k.ContainerResources(namespace)
	if err != nil {
		panic(err.Error())
	}

	capacity, err := k.ClusterCapacity()
	if err != nil {
		panic(err.Error())
	}

	nodes, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	rows := FormatResourceUsage(metricsclient, nodes, capacity, resources, sort, reverse, nodesOnly)

	if csv {
		prefix := "kube-resource-usage"
		if namespace == "" {
			prefix += "-all"
		} else {
			prefix += fmt.Sprintf("-%s", namespace)
		}

		filename := ExportCSV(prefix, rows)
		fmt.Printf("Exported %d rows to %s\n", len(rows), filename)
	} else {
		PrintResourceUsage(rows)
	}
}
