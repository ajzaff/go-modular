package control

import (
	"context"
	"math"

	"github.com/ajzaff/go-modular"
)

func Mul(ctx context.Context, a, b <-chan V, quit <-chan struct{}) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
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

func Sine(ctx context.Context, vs <-chan V) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
	go func() {
		for v := range vs {
			ch <- V(math.Sin(2 * math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}

func Sawtooth(ctx context.Context, vs <-chan V) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
	go func() {
		for v := range vs {
			ch <- V(2 / math.Pi * math.Atan(math.Tan(math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Triangle(ctx context.Context, vs <-chan V) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
	go func() {
		for v := range vs {
			ch <- V(2 / math.Pi * math.Asin(math.Sin(2*math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Sinc(ctx context.Context, vs <-chan V) <-chan V {
	ch := make(chan V, modular.BufferSize(ctx))
	go func() {
		for v := range vs {
			ch <- V(math.Sin(math.Pi*float64(v)) / (math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}
