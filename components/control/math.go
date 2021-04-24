package control

import (
	"context"
	"math"

	"github.com/ajzaff/go-modular"
)

func Mul(ctx context.Context, a, b CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
	loop:
		for {
			ch <- <-a * <-b
			select {
			case <-ctx.Done():
				break loop
			default:
			}
		}
		close(ch)
	}()
	return ch
}

func Sine(ctx context.Context, cv CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for v := range cv {
			ch <- modular.V(math.Sin(2 * math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}

func Sawtooth(ctx context.Context, cv CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for v := range cv {
			ch <- modular.V(2 / math.Pi * math.Atan(math.Tan(math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Triangle(ctx context.Context, cv CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for v := range cv {
			ch <- modular.V(2 / math.Pi * math.Asin(math.Sin(2*math.Pi*float64(v))))
		}
		close(ch)
	}()
	return ch
}

func Sinc(ctx context.Context, cv CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for v := range cv {
			ch <- modular.V(math.Sin(math.Pi*float64(v)) / (math.Pi * float64(v)))
		}
		close(ch)
	}()
	return ch
}
