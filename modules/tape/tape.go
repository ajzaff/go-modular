package modular

import (
	"fmt"
	"io"
)

type Buffer struct {
	buf []float32
	p   int64
}

func (t *Buffer) Pos() int64 {
	return t.p
}

func (t *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekCurrent:
		t.p += offset
	case io.SeekEnd:
		t.p = int64(len(t.buf)) - t.p
	case io.SeekStart:
		t.p = offset
	default:
		return 0, fmt.Errorf("tape.Buffer.Seek: invalid whence")
	}
	if t.p < 0 {
		t.p = 0
	} else if t.p > int64(len(t.buf)) {
		t.p = int64(len(t.buf))
	}
	return t.p, nil
}
