package vca

import (
	"io"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modio"
)

// VCA is a simple voltage controlled amplifier.
//
// CV in is the audio signal.
// CV a is the amplitude voltage usually sourced from an ADSR envelope generator.
type VCA struct {
	a   modio.Buffer
	in  modio.Buffer
	buf []modular.V
}

func (a *VCA) BlockSize() int { return 512 }

func (a *VCA) Write(vs []modular.V) (n int, err error) {
	return a.in.Write(vs)
}

func (a *VCA) Read(vs []modular.V) (n int, err error) {
	if a.buf == nil {
		a.buf = make([]modular.V, 512)
	}
	l := len(vs)
	if p := len(a.buf); p < l {
		l = p
	}
	if p := a.a.Len(); p < l {
		l = p
	}
	if p := a.in.Len(); p < l {
		l = p
	}
	buf := a.buf
	if l < len(a.buf) {
		buf = buf[:l]
	}
	_, err = a.a.Read(buf)
	n, err1 := a.in.Read(vs[:l])
	for i, v := range vs[:n] {
		vs[i] = a.buf[i] * v
	}
	if err == nil {
		err = err1
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return n, err
}

func (a *VCA) A() modular.Writer {
	return &vcaA{a: &a.a}
}

type vcaA struct {
	a *modio.Buffer
}

func (a *vcaA) Write(vs []modular.V) (n int, err error) {
	return a.a.Write(vs)
}
