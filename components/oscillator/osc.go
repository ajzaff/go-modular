package osc

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
)

const mid = 440

func audio(ctx *modular.Context, cin <-chan control.V) <-chan modular.V {
	ch := make(chan modular.V, ctx.BufferSize)
	go func() {
		for v := range cin {
			ch <- modular.V(v)
		}
		close(ch)
	}()
	return ch
}

func Sine(ctx *modular.Context, a <-chan control.V, freq <-chan control.V, quit <-chan struct{}) <-chan modular.V {
	var i uint64
	return audio(ctx, control.Mul(ctx, a, control.Sine(ctx, control.Func(ctx, func() (v control.V) {
		v.Store(float64(mid+<-freq) * float64(i) / float64(ctx.SampleRate))
		i++
		return
	})), quit))
}

func Sawtooth(ctx *modular.Context, freq <-chan control.V) <-chan modular.V {
	var i uint64
	return audio(ctx, control.Sawtooth(ctx, control.Func(ctx, func() (v control.V) {
		v.Store(float64(mid+<-freq) * float64(i) / float64(ctx.SampleRate))
		i++
		return
	})))
}

func Triangle(ctx *modular.Context, freq <-chan control.V) <-chan modular.V {
	var i uint64
	return audio(ctx, control.Triangle(ctx, control.Func(ctx, func() (v control.V) {
		v.Store(float64(mid+<-freq) * float64(i) / float64(ctx.SampleRate))
		i++
		return
	})))
}
