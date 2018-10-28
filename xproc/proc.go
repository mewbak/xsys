package xproc

// https://github.com/fastly/go-utils/blob/master/executable/executable.go
// https://github.com/crquan/coremem/blob/master/coremem.go
// https://github.com/janimo/memchart/blob/master/memchart.go

import (
	"github.com/kardianos/osext"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/process"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"github.com/smcduck/xsys/xcmd"
	"github.com/smcduck/xsys/xfs"
)

// on linux
// 32768 by default, you can read the value on your system in /proc/sys/kernel/pid_max,
// and you can set the value higher (up to 32768 for 32 bit systems or 4194304 for 64 bit) with:
// echo 4194303 > /proc/sys/kernel/pid_max
// on windows
// process id is a DWORD value, so max id value is int32 max value 4294967295
type ProcId int32

const (
	InvalidProcId ProcId = -1
)

type ProcInfo struct {
	Name           string  // short filename
	Filename       string  // full file path
	Cmdline        string  // run command line
	Param          string  // run param
	MemUsedBytes   uint64  // FIXME: unsupported for now
	CpuUsedPercent float64 // FIXME: unsupported for now
}

func GetAllPids() ([]ProcId, error) {
	pids, err := process.Pids()
	var result []ProcId
	for _, pid := range pids {
		result = append(result, ProcId(pid))
	}
	return result, err
}

func GetPidOfMyself() ProcId {
	pid := os.Getpid()
	return ProcId(pid)
}

func GetPidByProcFullFilename(filename string) ([]ProcId, error) {
	return nil, nil
}

func GetPidByProcName(procName string) ([]ProcId, error) {
	return nil, nil
}

// 替换自己的二进制文件
func ReplaceMySelfFile(newfilepath string) error {
	return nil
}

// 还可以参考：https://github.com/rcrowley/goagain/blob/master/goagain.go#L77
// Restart current process, with same parameters.
func RestartMyself() error {
	argv0, err := GetMyFullFilename()
	if err != nil {
		return err
	}
	files := make([]*os.File, syscall.Stderr+1)
	files[syscall.Stdin] = os.Stdin
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr
	wd, err := os.Getwd()
	if nil != err {
		return err
	}
	_, err = os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
	os.Exit(0)
	return err
}

// Notice:
// _, b, _, _ := runtime.Caller(0)
// return filepath.Dir(b)
// this is wrong
//
// Get process file folder, not working folder
func GetMyFolder() (string, error) {
	p, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

func GetMyFullFilename() (string, error) {
	return osext.Executable()
}

// TODO: replace GetMyFolder GetMyFullFilename
func SelfPath() (fullname, shortname, dir string, err error) {
	fullname, err = osext.Executable()
	if err != nil {
		return "", "", "", err
	}
	dir = filepath.Dir(fullname)
	shortname = xfs.FullFilenameToShort(fullname)
	return fullname, shortname, dir, err
}

func Terminate(pid ProcId) error {
	if pid < 0 {
		return errors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}
	return proc.Terminate()
}

// 隐藏进程,通过ps等命令看不到进程信息, 一般是通过把pid改为0做到的,windows下应该有api
func Hide(pid ProcId) error {
	if pid < 0 {
		return errors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	return nil
}

func Show(pid ProcId) error {
	if pid < 0 {
		return errors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}
	return nil
}

// remove run param from command line, get the filename
// samples:
// /usr/bin/spindump -> /usr/bin/spindump
// /usr/bin/spindump -usr/bin ->/usr/bin/spindump
// /usr/bin/spindump usr -> /usr/bin/spindump
func cmdlineToFilename(cmdline string) string {
	cr := []rune(cmdline)
	for i := 0; i < len(cmdline); i++ {
		if cmdline[i] == ' ' {
			exe := string(cr[0:i])
			_, err := os.Stat(exe)
			if err == nil {
				return exe
			}
		}
	}
	return cmdline
}

func GetProcInfo(pid ProcId) (*ProcInfo, error) {
	if pid < 0 {
		return nil, errors.New("invalid process id " + strconv.FormatInt(int64(pid), 10))
	}

	var pi ProcInfo
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return &pi, err
	}
	pi.Name, err = proc.Name()
	if err != nil {
		pi.Name = ""
	}
	cmdline, err := proc.Cmdline()
	pi.Cmdline = cmdline
	if runtime.GOOS == "darwin" {
		pi.Filename = cmdlineToFilename(cmdline)
		// mac下用cmdlineToFilename有个bug，就是使用./的方式执行的程序，在mac下得到的command不包含完整路径
		// 执行./goecho的时候，查询到的完整路径Filename是"./goecho"，不是完整路径
		// 所以使用了下述解决办法
		if len(pi.Filename) == 0 || strings.Index(pi.Filename, "./") == 0 {
			cmdstr := "lsof -p " + strconv.FormatInt(int64(pid), 10) + " | grep 'txt.*" + pi.Name + "'"
			result := xcmd.ExecWait(cmdstr, false)
			result = strings.Trim(result, "\r")
			result = strings.Trim(result, "\n")
			if len(result) > 0 {
				begin := strings.LastIndex(result, " /")
				if begin >= 0 {
					pi.Filename = result[begin+1:]
				}
			}
		}
	} else {
		pi.Filename, err = proc.Exe()
	}
	pi.Param = strings.Replace(cmdline, pi.Filename, "", 1)
	pi.Param = strings.TrimSpace(pi.Param)
	pi.CpuUsedPercent = -1 // Unsupported for now
	mi, err := proc.MemoryInfo()
	if err != nil {
		pi.MemUsedBytes = mi.Swap
	} else {
		pi.MemUsedBytes = 0
	}
	return &pi, nil
}
