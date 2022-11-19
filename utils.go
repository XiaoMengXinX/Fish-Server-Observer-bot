package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func GetCPUPercents() string {
	stats, _ := cpu.Percent(time.Second, false)
	return fmt.Sprintf("CPU: %.2f%%", stats[0])
}

func GetCPUCoresPercents() string {
	stats, _ := cpu.Percent(time.Second, true)
	var cores string
	for i, core := range stats {
		cores += fmt.Sprintf("Core %d: %.2f%%", i, core)
		if i != len(stats)-1 {
			cores += "\n"
		}
	}
	return cores
}

func GetMemStats() string {
	stats, _ := mem.VirtualMemory()
	var memory string
	memory += fmt.Sprintf("Memory: %.2f%% (%s/%s)", stats.UsedPercent, byteCountIEC(int64(stats.Used)), byteCountIEC(int64(stats.Total)))
	if stats.SwapTotal > 0 {
		swapUsed := stats.SwapTotal - stats.SwapFree
		swapPercent := float64(swapUsed) / float64(stats.SwapTotal) * 100
		memory += fmt.Sprintf("\nSwap: %.2f%% (%s/%s)", swapPercent, byteCountIEC(int64(swapUsed)), byteCountIEC(int64(stats.SwapTotal)))
	}
	return memory
}

func GetRootUsage() string {
	stats, _ := disk.Usage("/")
	return fmt.Sprintf("Usage of /: %.2f%% (%s/%s)", stats.UsedPercent, byteCountIEC(int64(stats.Used)), byteCountIEC(int64(stats.Total)))
}

func GetPartsStats() string {
	parts, _ := disk.Partitions(true)
	var usage string
	for i, part := range parts {
		stats, _ := disk.Usage(part.Mountpoint)
		usage += fmt.Sprintf("%s:\t%.2f%% (%s/%s)", part.Mountpoint, stats.UsedPercent, byteCountIEC(int64(stats.Used)), byteCountIEC(int64(stats.Total)))
		if i != len(parts)-1 {
			usage += "\n"
		}
	}
	return usage
}

func GetNetworkAllStats() string {
	sleepTime := 500
	stats, _ := net.IOCounters(false)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	stats2, _ := net.IOCounters(false)
	in := float64(stats2[0].BytesRecv-stats[0].BytesRecv) / (float64(sleepTime) / 1000)
	out := float64(stats2[0].BytesSent-stats[0].BytesSent) / (float64(sleepTime) / 1000)
	netStats := fmt.Sprintf("Network: %s/s in  %s/s out", byteCountIEC(int64(in)), byteCountIEC(int64(out)))
	return netStats
}

func GetNetworkStats() string {
	sleepTime := 500
	stats, _ := net.IOCounters(true)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	stats2, _ := net.IOCounters(true)
	var netStats string
	for i, stat := range stats {
		in := float64(stats2[i].BytesRecv-stat.BytesRecv) / (float64(sleepTime) / 1000)
		out := float64(stats2[i].BytesSent-stat.BytesSent) / (float64(sleepTime) / 1000)
		netStats += fmt.Sprintf("%s:\t%s/s in  %s/s out", stat.Name, byteCountIEC(int64(in)), byteCountIEC(int64(out)))
		if i != len(stats)-1 {
			netStats += "\n"
		}
	}
	return netStats
}

func GetHostInfo() string {
	info, _ := host.Info()
	os := firstUpper(info.OS)
	platform := firstUpper(info.Platform)
	var hostInfo string
	hostInfo += fmt.Sprintf("OS: %s %s %s (%s)\n", os, platform, info.PlatformVersion, info.KernelArch)
	hostInfo += fmt.Sprintf("Kernel: %s\n", info.KernelVersion)
	uptime := time.Duration(info.Uptime) * time.Second
	hostInfo += fmt.Sprintf("Uptime: %s", formatUptime(uptime))
	return hostInfo
}
