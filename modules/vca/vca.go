package vca

import (
	"github.com/ajzaff/go-modular"
)

// VCA is a simple voltage controlled amplifier.
//
// CV in is the audio signal.
// CV a is the amplitude voltage usually sourced from an ADSR envelope generator.
type VCA struct {
	a  []modular.V
	ap int

	in  []modular.V
	inp int
}

func (a *VCA) BlockSize() int { return 512 }

func (a *VCA) Write(vs []modular.V) (n int, err error) {
	if a.in == nil {
		a.in = make([]modular.V, 512)
		a.inp = 0
	}
	n = copy(a.in[a.inp:], vs)
	a.inp += n
	return n, nil
}

func (a *VCA) Read(vs []modular.V) (n int, err error) {
	for i := range vs {
		var v, va modular.V
		if i < a.inp {
			v = a.in[i]
		}
		if i < a.ap {
			va = a.a[i]
		}
		vs[i] = va * v
	}
	a.inp = 0
	a.ap = 0
	return n, nil
}

func (a *VCA) A() modular.Writer {
	return &vcaA{a: &a.a, ap: &a.ap}
}

type vcaA struct {
	a  *[]modular.V
	ap *int
}

func (a *vcaA) Write(vs []modular.V) (n int, err error) {
	if *a.a == nil {
		*a.a = make([]modular.V, 512)
		*a.ap = 0
	}
	n = copy((*a.a)[*a.ap:], vs)
	*a.ap += n
	return n, nil
}
