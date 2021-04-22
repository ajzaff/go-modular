package control

import (
	"context"

	"github.com/ajzaff/go-modular"
)

// V is a control voltage constant.
//
// The alias can be used to differentiate where a V is expected.
type V modular.V

// Store the float val into v.
func (v *V) Store(val float64) { *v = V(val) }

// Voltage returns a constant voltage source from v.
//
// Same as calling Func with a constant yielding fn.
func Voltage(ctx context.Context, cv V) <-chan V {
	return Func(ctx, func() V { return cv })
}

// Func returns a variable voltage source from evaluating fn.
func Func(ctx context.Context, fn func() V) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
	go func() {
		for {
			ch <- fn()
		}
	}()
	return ch
}
