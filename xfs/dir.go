package xfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func MakeDir(dir string) error {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	return os.MkdirAll(dir, 0777)
}

func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}

// Remove all content under dir, but keep dir folder
func CleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func ListDir(dir string) (dirs []string, files []string, err error) {
	if err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				dirs = append(dirs, path)
			} else {
				files = append(files, path)
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}
	return dirs, files, nil
}

func ListDirContains(dir, contains string) (dirs []string, files []string, err error) {
	if err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if strings.Contains(path, contains) {
				if f.IsDir() {
					dirs = append(dirs, path)
				} else {
					files = append(files, path)
				}
			}
			return nil
		}); err != nil {
		return nil, nil, err
	}
	return dirs, files, nil
}

func DirSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
