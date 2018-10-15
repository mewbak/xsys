package xio

import (
	"io"
	"time"
)

type SetDeadlineCallback func(t time.Time) error

// Forked from standard library io.Copy
func CopyTimeout(dst io.Writer, dstWriteCb SetDeadlineCallback, src io.Reader, srcReadCb SetDeadlineCallback, timeout time.Duration) (written int64, err error) {
	buf := make([]byte, 32*1024)
	var nr, nw int
	var er, ew error

	/*
		// If the reader has a WriteTo method, use it to do the copy.
		// Avoids an allocation and a copy.
		if wt, ok := src.(WriterTo); ok {
			return wt.WriteTo(dst)
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(ReaderFrom); ok {
			return rt.ReadFrom(src)
		}
		if buf == nil {
			buf = make([]byte, 32*1024)
		}
	*/

	for {
		if timeout > 0 {
			srcReadCb(time.Now().Add(timeout))
		}
		nr, er = src.Read(buf)
		if nr > 0 {
			if timeout > 0 {
				dstWriteCb(time.Now().Add(timeout))
			}
			nw, ew = dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
