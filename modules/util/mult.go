package util

import (
	"context"

	"github.com/ajzaff/go-modular"
)

// Mult copies the input audio signal n times.
//
// Note: all outputs must be utilized to avoid deadlocking.
func Mult(ctx context.Context, n int, in <-chan modular.V) []chan modular.V {
	if n <= 0 {
		panic("util.Mult: mult with <= 0 outputs")
	}
	out := make([]chan modular.V, n)
	for i := range out {
		out[i] = make(chan modular.V, modular.BufferSize(ctx))
	}
	go func() {
		for v := range in {
			for _, ch := range out {
				ch <- v
			}
		}
		for _, ch := range out {
			close(ch)
		}
	}()
	return out
}
