package xenvvar

// Notice, there is a library of homeDir: "github.com/mitchellh/go-homedir"

import (
	"bytes"
	"github.com/getlantern/appdir"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"github.com/MrMcDuck/xsys/xfs"
	"github.com/MrMcDuck/xapputil/xerror"
	"github.com/MrMcDuck/xsys/xcmd"
)

func SysRootDir() (string, error) {
	if "windows" == runtime.GOOS {
		return os.Getenv("SYSTEMDRIVE") + "\\", nil
	} else {
		return "/", nil
	}
}

func GetSharedUserDir() (string, error) {
	root, err := SysRootDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return root + "Users" + xfs.DirSlash() + "Public", nil
	}
	if runtime.GOOS == "darwin" {
		return root + "Users" + xfs.DirSlash() + "Shared", nil
	}
	return "", xerror.New("GetSharedUserDir() doesn't support %s", runtime.GOOS)
}

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func GetHomeDir() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support
	if "windows" == runtime.GOOS {
		return homeDirWindows()
	}
	// Unix-like system, so just assume Unix
	return homeDirUnix()
}

func homeDirUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeDirWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func GetGoRoot() (string, error) {
	// go env GOROOT可以正常返回的情况下，下述方法返回值不一定对，可能为空，所以部采用
	//os.Getenv("GOPATH"), nil
	goroot := xcmd.ExecWait("go env GOROOT", false) // 目前返回的结果最后会多一个换行符，是ExecWait的bug还是本来应该如此，尚不明确
	goroot = strings.Trim(goroot, "\r")
	goroot = strings.Trim(goroot, "\n")
	return goroot, nil
}

func GetGoPath() (string, error) {
	gopath := xcmd.ExecWait("go env GOPATH", false)
	gopath = strings.Trim(gopath, "\r")
	gopath = strings.Trim(gopath, "\n")
	return gopath, nil
}

func GetSystemLang() {
}

func SetSystemLang() {
}

func GetSocketLimit() uint32 {
	return 0
}

func SetSocketLimit(size uint32) {
}

func GetTimeWaitReuse() bool {
	return false
}

// echo "1" > /proc/sys/net/ipv4/tcp_tw_reuse
func SetTimeWaitReuse(onoff bool) {
}

func GetTimeWaitRecycle() bool {
	return false
}

// echo "1" > /proc/sys/net/ipv4/tcp_tw_recycle
func SetTimeWaitRecycle(onoff bool) {
}

func GetAppBinFolder(name string) string {
	return appdir.General(name)
}

// Windows:
// Linux:
// MacOS:  ~/Library/Logs/<App>
func GetAppLogFolder(name string) string {
	return appdir.Logs(name)
}