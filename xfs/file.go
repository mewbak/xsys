package xfs

import (
	"io"
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/pkg/errors"
)

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func FilenameToReader(filename string) (io.Reader, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return io.Reader(file), file, nil
}

func FilenameToWriter(filename string) (io.Writer, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return io.Writer(file), file, nil
}

func FileToBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func FileToJson(filename string, ptrJsonStruct interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, ptrJsonStruct)

	// Another implement
	/*
	file, err := os.Open(filename) // For read access.
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(ptrJsonStruct)*/
}

func JsonToFile(jsonStruct interface{}, indent bool, filename string) error {
	if jsonStruct == nil {
		return errors.New("Null input jsonStruct")
	}
	b := []byte{}
	err := error(nil)
	if indent {
		b, err = json.MarshalIndent(jsonStruct, "", "\t")
	} else {
		b, err = json.Marshal(jsonStruct)
	}
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	f.Write(b)
	f.Close()
	return nil
}

func BytesToFile(data []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	total := len(data)
	done := 0
	for {
		if done >= total {
			break
		}
		n, err := f.Write(data)
		if err != nil {
			return err
		}
		done += n
	}
	return nil
}

func StringToFile(data string, filename string) error {
	return BytesToFile([]byte(data), filename)
}

func AppendStringToFile(data string, filename string) error {
	return AppendBytesToFile([]byte(data), filename)
}

func AppendBytesToFile(data []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}
	return nil
}

func GetFileSize(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
