package xwnd

import (
	"github.com/go-vgo/robotgo"
	"github.com/smcduck/xsys/xfs"
	"github.com/smcduck/xsys/xproc"
	"github.com/pkg/errors"
	"strings"
)

type WndInfo struct {
	Title    string
	Pid      uint64
	Filename string
	App      string
}

func checkAppByWndInfo(wi WndInfo) string {
	if len(wi.Filename) == 0 {
		return "unknown"
	}
	fn := xfs.FullFilenameToShort(wi.Filename)
	fn = strings.ToLower(fn)
	fn = strings.Replace(fn, "google ", "", -1)
	return fn
}

// https://play.golang.org/p/YfGDtIuuBw // windows only
// https://gist.github.com/obonyojimmy/d6b263212a011ac7682ac738b7fb4c70 // windows only
// http://codegists.com/snippet/go/active-windowgo_obonyojimmy_go // windows only
// 此函数涉及安全性API, 需要在"安全与隐私 -> 辅助功能"中将终端或者代码父进程加入白名单
func GetActiveWindowInfo() (WndInfo, error) {
	var wi WndInfo
	wi.Title = robotgo.GetTitle()
	if wi.Title == "IsValid failed." {
		wi.Title = ""
	}
	wi.Pid = uint64(robotgo.GetPID())
	if wi.Pid < 0 {
		return wi, errors.New("Get PID fail")
	}
	pi, err := xproc.GetProcInfo(xproc.ProcId(wi.Pid))
	if err != nil {
		return wi, err
	}
	wi.Filename = pi.Filename
	wi.App = checkAppByWndInfo(wi)
	return wi, nil
}
