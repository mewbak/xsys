package xhdd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/disk"
	"github.com/MrMcDuck/xdsa/xvolume"
)

type PartitionInfo struct {
	FilesystemType string
	AvailableBytes uint64
	FreeBytes      uint64
	TotalBytes     uint64
}

func (pi PartitionInfo) String() string {
	total, err := xvolume.NewFromByteSize(float64(pi.TotalBytes))
	if err != nil {
		return err.Error()
	}
	free, err := xvolume.NewFromByteSize(float64(pi.FreeBytes))
	if err != nil {
		return err.Error()
	}
	available, err := xvolume.NewFromByteSize(float64(pi.AvailableBytes))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintln(
		"FilesystemType:", pi.FilesystemType,
		"Total:", total.String(),
		"Free:", free.String(),
		"Available:", available.String(),
	)
}

// Get returns disk usage info and error if any.
func GetPartitionInfo(partitionPath string) (*PartitionInfo, error) {
	var pi PartitionInfo
	du, err := disk.Usage(partitionPath)
	if err != nil {
		return nil, err
	}
	pi.AvailableBytes = du.Total
	pi.FreeBytes = du.Total
	pi.TotalBytes = du.Total
	pi.FilesystemType = du.Fstype
	return &pi, nil
}

// returns partition path, same as mountpoint in Unix, or logical drive in Windows
func ListPartitions() (partitionPath []string, err error) {
	ps, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, errors.New("Get disk partitions error")
	}

	var result []string
	for _, item := range ps {
		result = append(result, item.Mountpoint)
	}
	return result, nil
}

// https://github.com/cydev/du
// https://github.com/ricochet2200/go-disk-usage
// https://gist.github.com/lunny/9828326
// http://wendal.net/2012/1224.html
// https://github.com/lxn/win
// https://github.com/AllenDang/w32

/* Custom implement for Windows
func getPartitionUsage(partitionPath string) (PartitionUsage, error) {
	u := PartitionUsage{}
	h, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return u, err
	}
	c, err := h.FindProc("GetDiskFreeSpaceExW")
	if err != nil {
		return u, err
	}
	_, _, err = c.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(partitionPath))),
		uintptr(unsafe.Pointer(&u.Free)),
		uintptr(unsafe.Pointer(&u.Total)),
		uintptr(unsafe.Pointer(&u.Available)))

	return u, err
}


// http://stackoverflow.com/questions/23128148/how-can-i-get-a-listing-of-all-drives-on-windows-using-golang
// 正确答案貌似更加简单
func listPartitions() ([]string, error) {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return nil, err
	}
	getLogicalDrivesHandle, err := syscall.GetProcAddress(kernel32, "GetLogicalDrives")
	if err != nil {
		return nil, err
	}
	if ret, _, callErr := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); callErr != 0 {
		return nil, callErr
	} else { // bitsToPartitions
		partitions := make([]string, 0)
		allPartitions := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		for i := 0; i < 26; i++ {
			if bitMap == 1 {
				partitions = append(partitions, allPartitions[i] + ":")
			}
			bitMap >>= 1
		}
		return partitions, nil
	}
}
*/

/* Custom implement for Unix
func getPartitionUsage(partitionPath string) (PartitionUsage, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(partitionPath, &stat)
	i := PartitionUsage {
		Total:     int64(stat.Bsize) * int64(stat.Blocks),
		Available: int64(stat.Bsize) * int64(stat.Bavail),
		Free:      int64(stat.Bsize) * int64(stat.Bfree),
	}
	return i, err
}
*/
