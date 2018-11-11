package xcmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// sudo 下执行ExecWait或者ExecNoWait接口时会有问题
// 频繁循环调用同一个命令可能报错"fork/exec /usr/local/bin/ffmpeg: bad file descriptor"，也可能不是频繁调用的锅

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

func ExecWait(cmdStr string, screenPrint bool) string {
	var cmd *exec.Cmd

	cmdStr = strings.TrimSpace(cmdStr)
	// FIXME 参数当中带myapp -d 'a b c'
	/*if strings.Count(cmdStr, " ") > 0 { // example: ping -c 3 baidu.com
		strs := strings.Split(cmdStr, " ")
		// slice打散语法糖, 将数组对应到可变参数列表上
		cmd = exec.Command(strs[0], strs[1:]...)
	} else { // example: ffmpeg
		cmd = exec.Command(cmdStr)
	}*/
	cmd = exec.Command(cmdStr)

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

type Cmder struct {
	cmd *exec.Cmd
}

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
}

func ExecNowait(cmdStr string, screenPrint bool) *Cmder {
	var cmder Cmder

	cmdStr = strings.TrimSpace(cmdStr)
	if strings.Count(cmdStr, " ") > 0 { // example: ping -c 3 baidu.com
		strs := strings.Split(cmdStr, " ")
		// slice打散语法糖, 将数组对应到可变参数列表上
		cmder.cmd = exec.Command(strs[0], strs[1:]...)
	} else { // example: ffmpeg
		cmder.cmd = exec.Command(cmdStr)
	}

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
