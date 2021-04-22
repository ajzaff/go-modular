package control

import (
	"math"

	"github.com/ajzaff/go-modular"
)

func Mul(ctx *modular.Context, a, b <-chan V, quit <-chan struct{}) <-chan V {
	ch := make(chan V, ctx.BufferSize)
	go func() {
	loop:
		for {
			ch <- <-a * <-b
			select {
			case <-quit:
				break loop
			default:
			}
		}
		close(ch)
	}()
	return ch
}

func Sine(ctx *modular.Context, vs <-chan V) <-chan V {
	ch := make(chan V, ctx.BufferSize)
	go func() {
		for v := range vs {
			ch <- V(math.Sin(2 * math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}

func Sawtooth(ctx *modular.Context, vs <-chan V) <-chan V {
	ch := make(chan V, ctx.BufferSize)
	go func() {
		for v := range vs {
			ch <- V(2 / math.Pi * math.Atan(math.Tan(math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Triangle(ctx *modular.Context, vs <-chan V) <-chan V {
	ch := make(chan V, ctx.BufferSize)
	go func() {
		for v := range vs {
			ch <- V(2 / math.Pi * math.Asin(math.Sin(2*math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Sinc(ctx *modular.Context, vs <-chan V) <-chan V {
	ch := make(chan V, ctx.BufferSize)
	go func() {
		for v := range vs {
			ch <- V(math.Sin(math.Pi*float64(v)) / (math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}
