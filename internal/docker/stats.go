package docker

import (
	"context"
	"encoding/json"
)

type Stat struct {
	ID     string
	CPU    float64
	Memory uint64
}

type cpuStats struct {
	CPUUsage struct {
		TotalUsage uint64 `json:"total_usage"`
	} `json:"cpu_usage"`
	SystemUsage uint64 `json:"system_cpu_usage"`
	OnlineCPUs  uint32 `json:"online_cpus"`
}

type apiStats struct {
	CPU    cpuStats `json:"cpu_stats"`
	PreCPU cpuStats `json:"precpu_stats"`
	Memory struct {
		Usage uint64 `json:"usage"`
	} `json:"memory_stats"`
}

func (c *Client) Stats(ctx context.Context, id string, out chan<- Stat) error {
	body, err := c.get(ctx, "/containers/"+id+"/stats?stream=1")
	if err != nil {
		return err
	}
	defer body.Close()

	dec := json.NewDecoder(body)
	for {
		var s apiStats
		if err := dec.Decode(&s); err != nil {
			return err
		}

		cpuDelta := float64(s.CPU.CPUUsage.TotalUsage - s.PreCPU.CPUUsage.TotalUsage)
		sysDelta := float64(s.CPU.SystemUsage - s.PreCPU.SystemUsage)
		pct := 0.0
		if sysDelta > 0 && s.CPU.OnlineCPUs > 0 {
			pct = cpuDelta / sysDelta * float64(s.CPU.OnlineCPUs) * 100
		}

		select {
		case out <- Stat{ID: id, CPU: pct, Memory: s.Memory.Usage}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
