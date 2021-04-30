package control

import (
	"github.com/ajzaff/go-modular"
)

// Voltage returns a CV from a singluar value v.
//
// Equivalent to calling Func with a constant-yielding func.
func Voltage(v float64) modular.Reader {
	return Func(func() modular.V { return modular.V(v) })
}

// Func returns a variable voltage source from evaluating fn.
func Func(fn func() modular.V) modular.Reader {
	return controlFunc(fn)
}

type controlFunc func() modular.V

func (f controlFunc) Read(vs []modular.V) (n int, err error) {
	for i := range vs {
		vs[i] = f()
	}
	return len(vs), nil
}
