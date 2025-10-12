package helpers

import (
	"expansion-gateway/dto/processes"
	"expansion-gateway/errors/processerror"
	"expansion-gateway/interfaces/errorinfo"

	"github.com/shirou/gopsutil/v4/process"
)

// returns the resource usage percent of a process in the RAM and the CPU
func GetResourceUsageOfProcess(pid int32) (*processes.ProcessData, errorinfo.GatewayError) {
	if p, err := process.NewProcess(pid); err == nil {
		cpuUsage, _ := p.CPUPercent()
		ramUsage, _ := p.MemoryPercent()

		return &processes.ProcessData{
			CPUusage: cpuUsage,
			RAMusage: ramUsage,
		}, nil
	}

	return nil, processerror.CreateProcessNotFoundError(
		"/helpers/get_resource_usage_of_process.go",
		12,
		pid,
	)
}

func GetResourceUsageOfProcessNoError(pid int32) *processes.ProcessData {
	if p, err := GetResourceUsageOfProcess(pid); err == nil {
		return p
	}

	return &processes.ProcessData{
		CPUusage: 0,
		RAMusage: 0,
	}
}
