package xcpu

import (
	//"github.com/klauspost/cpuid" // x86/x64 is supported only for now
	"github.com/shirou/gopsutil/cpu"
	"time"
	"runtime"
)

// Get unique serial number of CPU
func GetSerialNumber() (string, error) {
	return "Unsupported for now", nil
}

// 获取所有CPU的使用百分比，以数组返回
func GetAllUsedPercent(duration time.Duration) ([]float64, error) {
	return cpu.Percent(duration, true)
}

// 获取所有CPU的使用百分比，组合成总百分比后返回
func GetCombinedUsedPercent(duration time.Duration) (float64, error) {
	p, err := cpu.Percent(duration, false)
	if err != nil {
		return 0, err
	}
	return p[0], err
}

func GetCpuCount() int {
	return runtime.NumCPU()
}

