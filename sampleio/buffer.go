package sampleio

import (
	"errors"
	"io"

	"github.com/ajzaff/go-modular"
)

type Buffer struct {
	buf []modular.V
	off int
}

func NewBuffer(buf []modular.V) *Buffer {
	return &Buffer{buf: buf}
}

var errTooLarge = errors.New("sample.Buffer: too large")
var errNegativeRead = errors.New("sample.Buffer: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

func (b *Buffer) Samples() []modular.V { return b.buf[b.off:] }

func (b *Buffer) empty() bool { return len(b.buf) <= b.off }

func (b *Buffer) Len() int { return len(b.buf) - b.off }

func (b *Buffer) Cap() int { return cap(b.buf) }

func (b *Buffer) Truncate(n int) {
	if n == 0 {
		b.Reset()
		return
	}
	if n < 0 || n > b.Len() {
		panic("bytes.Buffer: truncation out of range")
	}
	b.buf = b.buf[:b.off+n]
}

func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

const smallBufferSize = 64

func (b *Buffer) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]modular.V, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(errTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

func (b *Buffer) Write(p []modular.V) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

const MinRead = 512

func (b *Buffer) ReadFrom(r Reader) (n int64, err error) {
	for {
		i := b.grow(MinRead)
		b.buf = b.buf[:i]
		m, e := r.Read(b.buf[i:cap(b.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		b.buf = b.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

func makeSlice(n int) []modular.V {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(errTooLarge)
		}
	}()
	return make([]modular.V, n)
}

func (b *Buffer) WriteTo(w Writer) (n int64, err error) {
	if nSamples := b.Len(); nSamples > 0 {
		m, e := w.Write(b.buf[b.off:])
		if m > nSamples {
			panic("sample.Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all samples should have been written, by definition of
		// Write method in io.Writer
		if m != nSamples {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	b.Reset()
	return n, nil
}

func (b *Buffer) WriteSample(c modular.V) error {
	m, ok := b.tryGrowByReslice(1)
	if !ok {
		m = b.grow(1)
	}
	b.buf[m] = c
	return nil
}

func (b *Buffer) Read(p []modular.V) (n int, err error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	return n, nil
}

func (b *Buffer) Next(n int) []modular.V {
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	return data
}

func (b *Buffer) ReadSample() (modular.V, error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		return 0, io.EOF
	}
	c := b.buf[b.off]
	b.off++
	return c, nil
}
