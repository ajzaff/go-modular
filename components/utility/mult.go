package utility

import (
	"github.com/ajzaff/go-modular"
)

func Mult(ctx *modular.Context, n int, in <-chan modular.V) []chan modular.V {
	if n <= 0 {
		panic("utility.Mult: mult with <= 0 outputs")
	}
	out := make([]chan modular.V, n)
	for i := range out {
		out[i] = make(chan modular.V, ctx.BufferSize)
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
