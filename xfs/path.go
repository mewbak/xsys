package xfs

import (
	"os"
	"path"
	"strings"
	"time"
	"github.com/MrMcDuck/xdsa/xstring"
)

type PathInfo struct {
	Exist        bool
	IsFolder     bool
	ModifiedTime time.Time
}

const (
	InvalidFilenameCharsWindows = "\"\\:/*?<>|“”"
)

func GetPathInfo(path string) (*PathInfo, error) {
	var pi PathInfo
	fi, err := os.Stat(path)
	if err == nil {
		pi.Exist = true
		pi.IsFolder = fi.IsDir()
		pi.ModifiedTime = fi.ModTime()
		return &pi, nil
	} else if err != nil && os.IsNotExist(err) {
		pi.Exist = false
		return &pi, nil
	} else {
		return &pi, err
	}
}

func FileExits(filename string) bool {
	pi, err := GetPathInfo(filename)
	if err != nil {
		return false
	}
	return !pi.IsFolder && pi.Exist
}

func FolderExits(folder string) bool {
	pi, err := GetPathInfo(folder)
	if err != nil {
		return false
	}
	return pi.IsFolder && pi.Exist
}

// Combine absolute path and relative path to get a new absolute path
func PathJoin(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}

func FullFilenameToShort(fullFilename string) string {
	fullFilename = strings.TrimSpace(fullFilename)
	if len(fullFilename) == 0 {
		return ""
	}
	idx := strings.LastIndex(fullFilename, "/")
	if idx < 0 {
		idx = strings.LastIndex(fullFilename, "\\")
	}
	if idx < 0 || idx == (len(fullFilename)-1) {
		return ""
	}
	result, err := xstring.SubstrAscii(fullFilename, idx+1, len(fullFilename)-1)
	if err != nil {
		return ""
	}
	return result
}

// Replace illegal chars for short filename / dir name, not multi-level directory
func RefactShortPathName(path string) string {
	var illegalChars = "/\\:*\"<>|"
	for _, c := range illegalChars {
		path = strings.Replace(path, string(c), "-", -1)
	}
	return path
}
