package control

import (
	"context"

	"github.com/ajzaff/go-modular"
)

// CV is a control voltage.
//
// The alias can be used to differentiate between audio and control voltages.
type CV <-chan modular.V

// Voltage returns a CV from a singluar value v.
//
// Equivalent to calling Func with a constant-yielding func.
func Voltage(ctx context.Context, v float64) CV {
	return Func(ctx, func() modular.V { return modular.V(v) })
}

// Func returns a variable voltage source from evaluating fn.
func Func(ctx context.Context, fn func() modular.V) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for {
			ch <- fn()
		}
	}()
	return ch
}
