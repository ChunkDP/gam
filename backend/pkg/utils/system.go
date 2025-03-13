package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// GetCPUUsage 获取CPU使用率
func GetCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return 0
	}
	return percent[0]
}

// GetMemoryUsage 获取内存使用率
func GetMemoryUsage() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return v.UsedPercent
}

// GetDiskUsage 获取磁盘使用率
func GetDiskUsage() float64 {
	parts, err := disk.Partitions(true)
	if err != nil {
		return 0
	}

	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err == nil {
			return usage.UsedPercent
		}
	}
	return 0
}

// GetNetworkIO 获取网络IO信息
func GetNetworkIO() string {
	io, err := net.IOCounters(false)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("BytesSent: %d, BytesRecv: %d", io[0].BytesSent, io[0].BytesRecv)
}

// GetLoadAverage 获取系统负载
func GetLoadAverage() string {
	// 使用 gopsutil 的 load 包
	avg, err := load.Avg()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.2f %.2f %.2f", avg.Load1, avg.Load5, avg.Load15)
}
