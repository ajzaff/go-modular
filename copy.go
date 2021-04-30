package modular

import (
	"errors"
	"io"
)

// copyBuffer copies values from src to dst using buf until an error is reached.
// Returns the number of bytes copied and the error returned.
//
// done allows the copy to be ended approximately on demand so long
// src or dst is not blocking.
func copyBuffer(dst Writer, src Reader, buf []V, done <-chan struct{}) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("modular.copyBuffer: empty buffer")
	}
	if wt, ok := src.(interface {
		WriteTo(w Writer) (n int64, err error)
	}); ok {
		return wt.WriteTo(dst)
	}
	if rf, ok := dst.(interface {
		ReadFrom(r Reader) (n int64, err error)
	}); ok {
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
			if size == 0 {
				size = 32 * 1024 // io.Copy default
			}
		}
		// TODO: add support for limit reader?
		buf = make([]V, size)
	}
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("modular.copyBuffer: invalid write")
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
