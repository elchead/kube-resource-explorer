package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
	"k8s.io/apimachinery/pkg/api/resource"
)

func cmp(r []*ResourceAllocation, field string, i, j int, reverse bool) bool {
	f1 := r[i].getField(field)
	f2 := r[j].getField(field)

	if q1, ok := f1.(resource.Quantity); ok {
		q2 := f2.(resource.Quantity)
		if reverse {
			return q1.Cmp(q2) < 0
		}
		return q1.Cmp(q2) > 0
	}

	if v1, ok := f1.(int64); ok {
		v2 := f2.(int64)
		if reverse {
			return v1 < v2
		}
		return v1 > v2

	}

	if s1, ok := f1.(string); ok {
		s2 := f2.(string)
		if reverse {
			return strings.Compare(s1, s2) > 0
		}
		return strings.Compare(s1, s2) < 0
	}

	panic("Unknown type")
}

func fmtPercent(p int64) string {
	return fmt.Sprintf("%d%%", p)
}

func (r *ResourceLister) Print(resourceUsage []*ResourceAllocation, field string, reverse bool) {

	sort.Slice(resourceUsage, func(i, j int) bool {
		return cmp(resourceUsage, field, i, j, reverse)
	})

	rows := []string{
		"Namespace | Name | CpuReq | CpuReq% | CpuLimit | CpuLimit% | MemReq | MemReq% | MemLimit | MemLimit%",
		"--------- | ---- | ------ | ------- | -------- | --------- | ------ | ------- | -------- | ---------",
	}

	for _, u := range resourceUsage {
		row := strings.Join([]string{
			u.Namespace,
			u.Name,
			u.CpuReq.String(),
			fmtPercent(u.PercentCpuReq),
			u.CpuLimit.String(),
			fmtPercent(u.PercentCpuLimit),
			u.MemReq.String(),
			fmtPercent(u.PercentMemoryReq),
			u.MemLimit.String(),
			fmtPercent(u.PercentMemoryLimit),
		}, "| ")
		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
}
