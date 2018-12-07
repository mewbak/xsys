package xcmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// sudo 下执行ExecWait或者ExecNoWait接口时会有问题
// 频繁循环调用同一个命令可能报错"fork/exec /usr/local/bin/ffmpeg: bad file descriptor"，也可能不是频繁调用的锅

type (
	Options struct { // TODO
		CpuLimit float64
		Timeout  time.Duration
		KillAfterExit bool
	}

	Cmder struct {
		cmd *exec.Cmd
	}
)

func ExecWaitByShell(cmdStr string, screenPrint bool) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		// WARNING: 如果用sudo执行主程序，或者使用sudo /bin/sh -c xxx 来执行命令，会报错command not found
		cmd = exec.Command("/bin/sh", "-c", cmdStr)
	}
	if screenPrint {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Println(err.Error())
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println(err.Error())
		}
		return ""
	} else {
		out, _ := cmd.CombinedOutput()
		return string(out)
	}
}

// will print screen
func ExecWaitPrintScreen(name string, arg ...string) error {
	var cmd *exec.Cmd

	cmd = exec.Command(name, arg...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// will not print screen
func ExecWaitReturn(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

/*
func ExecNowaitByShell(cmdStr string, screenPrint bool) *Cmder {
	var cmder Cmder
	if runtime.GOOS == "windows" {
		cmder.cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmder.cmd = exec.Command("/bin/sh", "-c", cmdStr)
	}
	if screenPrint {
		cmder.cmd.Stdout = os.Stdout
		cmder.cmd.Stderr = os.Stderr
		if err := cmder.cmd.Start(); err != nil {
			fmt.Println(err.Error())
		}
	}
	return &cmder
}*/

func ExecNowait(screenPrint bool, cmd string, arg ...string) *Cmder {
	var cmder Cmder
	cmder.cmd = exec.Command(cmd, arg...)

	fmt.Println("command line:")
	fmt.Println(cmd, strings.Join(append([]string{}, arg...), " "))
	if screenPrint {
		cmder.cmd.Stdout = os.Stdout
		cmder.cmd.Stderr = os.Stderr
		if err := cmder.cmd.Start(); err != nil {
			fmt.Println(err.Error())
		}
	}
	return &cmder
}

func (c *Cmder) GetPid() int32 {
	return int32(c.cmd.Process.Pid)
}

func (c *Cmder) Wait() string {
	c.cmd.Wait()
	out, _ := c.cmd.CombinedOutput()
	return string(out)
}

type Interacter struct {
}

func NewInteracter(cmd string) *Interacter {
	return nil
}

func (i *Interacter) Read() (output string, eixt bool) {
	return "", false
}

func (i *Interacter) Write(input string) (output string, eixt bool) {
	return "", false
}

func Interact(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 1)
	out := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself

	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()
	go func() {
		// here i try to grep sudo from stderr, but not work
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := e.Read(buf[t:])
			if err != nil && err.Error() != "EOF" {
				fmt.Println(err)
			}
			if s := string(buf[t:]); strings.Contains(s, "sudo") {
				fmt.Println("here")
				w.Write([]byte("123456\n"))
			} else {
			}
			t += n
		}
	}()
	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			if s := string(buf[t:]); strings.Contains(s, "[sudo]") {
				w.Write([]byte("ubuntu\n"))
			} else {
			}
			t += n
			if buf[t-2] == '$' { //assuming the $PS1 == 'sh-4.3$ '
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
