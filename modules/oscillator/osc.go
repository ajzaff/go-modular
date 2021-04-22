// Package osc provides a standard VCO and LFO.
package osc

import (
	"context"
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
)

// Range presents a pipe organ length setting.
// The zero value is LFO and higher values are
// octaves at 32hz doubling at each setting.
type Range int

const (
	RangeLow Range = iota // 1hz
	Range32               // 32hz
	Range16               // 64hz
	Range8                // 128hz
	Range4                // 256hz
	Range2                // 512hz
)

func (r Range) Hz(fine float64) float64 {
	return 16 * math.Pow(2, float64(r)+fine)
}

func osc(ctx context.Context, wave func() modular.V) <-chan modular.V {
	ch := make(chan modular.V)
	go func() {
		for {
			ch <- wave()
		}
	}()
	return ch
}

func Sine(ctx context.Context, r Range, fine control.V) <-chan modular.V {
	i := 0
	freq := r.Hz(float64(fine))
	return osc(ctx, func() (v modular.V) {
		v = modular.V(math.Sin(2 * math.Pi * freq * float64(i) / float64(modular.SampleRate(ctx))))
		i++
		return
	})
}

func Triangle(ctx context.Context, r Range, fine control.V) <-chan modular.V {
	i := 0
	freq := r.Hz(float64(fine))
	return osc(ctx, func() (v modular.V) {
		v = modular.V(2 / math.Pi * math.Asin(math.Sin(2*math.Pi*freq*float64(i)/float64(modular.SampleRate(ctx)))))
		i++
		return
	})
}

func Saw(ctx context.Context, r Range, fine control.V) <-chan modular.V {
	i := 0
	freq := r.Hz(float64(fine))
	return osc(ctx, func() (v modular.V) {
		v = modular.V(2 / math.Pi * math.Atan(math.Tan(math.Pi*freq*float64(i)/float64(modular.SampleRate(ctx)))))
		i++
		return
	})
}
