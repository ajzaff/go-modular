package sampleio

import (
	"errors"
	"io"

	"github.com/ajzaff/go-modular"
)

func Copy(ctx *modular.Context, dst Writer, src Reader) (n int64, err error) {
	return CopyBuffer(ctx, dst, src, nil)
}

// CopyBuffer values from src to dst using buf until an error is reached.
// Returns the number of bytes copied and the error returned.
func CopyBuffer(ctx *modular.Context, dst Writer, src Reader, buf []modular.V) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("sample.CopyBuffer: empty buffer")
	}
	return copyBuffer(ctx, dst, src, buf)
}

func copyBuffer(ctx *modular.Context, dst Writer, src Reader, buf []modular.V) (written int64, err error) {
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
			} else if v := ctx.BufferSize; v > 0 {
				size = v
			} else {
				size = 32 * 1024 // io.Copy default
			}
		}
		// TODO: add support for limit reader?
		buf = make([]modular.V, size)
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
