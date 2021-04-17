package sample

import (
	"errors"
	"io"
)

func Equal(a, b []Sample) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Repeat(a []Sample, count int) []Sample {
	if count == 0 {
		return []Sample{}
	}
	if count < 0 {
		panic("samples: negative Repeat count")
	} else if len(a)*count/count != len(a) {
		panic("samples: Repeat count causes overflow")
	}

	buf := make([]Sample, len(a)*count)
	bp := copy(buf, a)
	for bp < len(buf) {
		copy(buf[bp:], buf[:bp])
		bp *= 2
	}
	return buf
}

func Copy(dst Writer, src Reader) (n int64, err error) {
	return CopyBuffer(dst, src, nil)
}

// CopyBuffer values from src to dst using buf until an error is reached.
// Returns the number of bytes copied and the error returned.
func CopyBuffer(dst Writer, src Reader, buf []Sample) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("sample.CopyBuffer: empty buffer")
	}
	return copyBuffer(dst, src, buf)
}

func copyBuffer(dst Writer, src Reader, buf []Sample) (written int64, err error) {
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rf, ok := dst.(ReaderFrom); ok {
		return rf.ReadFrom(src)
	}
	if buf == nil {
		var size int
		if sp, ok := src.(Processor); ok {
			size = sp.BlockSize()
		}
		if size == 0 {
			if dp, ok := dst.(Processor); ok {
				size = dp.BlockSize()
			}
			size = 32 * 1024 // io.Copy default
		}
		// TODO: add support for limit reader?
		buf = make([]Sample, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write")
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
