// Package osc provides standard VCO and LFO waveforms.
package osc

import (
	"context"
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	"github.com/ajzaff/go-modular/components/midi"
)

// Polarity controls the polarity of waveform functions.
//
// Negative "inverts" the wave, while positive maintains
// the true shape. Conveniently, non-pole values can be
// used to set the amplitude.
type Polarity float64

const (
	Negative Polarity = 2*iota - 1 // inverted
	Positive                       // regular
)

// Range presents a pipe organ length setting.
// The zero value is LFO and higher values are
// octaves at 32hz doubling at each setting.
type Range int

const (
	RangeLo Range = iota // 1hz@0
	Range32              // 2hz@0
	Range16              // 4hz@0
	Range8               // 8hz@0
	Range4               // 16hz@0
	Range2               // 32hz@0
)

// Tone returns the tone frequency for the range and fine tuning.
func Tone(r Range, fine float64) float64 {
	return math.Pow(2, float64(r)+fine/12)
}

// Fine returns the fine tuning constant to tune the oscillators to t at Range8.
func Fine(t midi.Tuning) float64 {
	return 12*math.Log2(t.A4Hz()) - 105
}

// Sine outputs an sine audio wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Sine(ctx context.Context, a Polarity, r Range, fine float64, lin control.CV) <-chan modular.V {
	i := 0
	sampleRate := modular.SampleRate(ctx)
	return osc(ctx, a, func() (v modular.V) {
		length := float64(sampleRate) / Tone(r, fine+float64(<-lin))
		v.Store(math.Sin(2 * math.Pi * float64(i) / length))
		i++
		return
	})
}

// SineReader outputs sine audio wave from samples from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func SineReader(ctx *modular.Context, a Polarity, r Range, fine float64, lin control.CV) modular.Reader {
	rd := &sineReader{
		a:    a,
		r:    r,
		fine: fine,
		opts: ctx.Options,
	}
	go func() {
		for v := range lin {
			rd.v.Store(v)
		}
	}()
	return rd
}

type sineReader struct {
	a    Polarity
	r    Range
	fine float64
	i    uint64
	v    control.Atomic
	opts modular.Options
}

func (r *sineReader) Read(vs []modular.V) (n int, err error) {
	for i := range vs {
		length := float64(r.opts.SampleRate) / Tone(r.r, r.fine+float64(r.v.Load()))
		vs[i] = modular.V(math.Sin(2 * math.Pi * float64(r.i) / length))
		r.i++
	}
	return len(vs), nil
}

// Triangle outputs an triangle wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Triangle(ctx context.Context, a Polarity, r Range, fine float64, lin control.CV) <-chan modular.V {
	i := 0
	sampleRate := modular.SampleRate(ctx)
	return osc(ctx, a, func() (v modular.V) {
		length := float64(sampleRate) / Tone(r, fine+float64(<-lin))
		v.Store(2 / math.Pi * math.Asin(math.Sin(2*math.Pi*float64(i)/length)))
		i++
		return
	})
}

// Saw outputs an sawtooth wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Saw(ctx context.Context, a Polarity, r Range, fine float64, lin control.CV) <-chan modular.V {
	i := 0
	sampleRate := modular.SampleRate(ctx)
	return osc(ctx, a, func() (v modular.V) {
		length := float64(sampleRate) / Tone(r, fine+float64(<-lin))
		v.Store(2 / math.Pi * math.Atan(math.Tan(math.Pi*float64(i)/length)))
		i++
		return
	})
}

// Square outputs an square wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Square(ctx context.Context, a Polarity, r Range, fine float64, lin control.CV) <-chan modular.V {
	return Pulse(ctx, a, r, fine, control.Voltage(ctx, .5), lin)
}

// Pulse outputs an pulse wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
// Pulse width w are in the range 0 to 1.
func Pulse(ctx context.Context, a Polarity, r Range, fine float64, w, lin control.CV) <-chan modular.V {
	i := 0
	sampleRate := modular.SampleRate(ctx)
	return osc(ctx, a, func() (v modular.V) {
		length := float64(sampleRate) / Tone(r, fine+float64(<-lin))
		if math.Mod(float64(i)/length, 2) < 2*float64(<-w) {
			v.Store(1)
		} else {
			v.Store(0)
		}
		i++
		return
	})
}

func osc(ctx context.Context, a Polarity, wave func() modular.V) <-chan modular.V {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for {
			ch <- modular.V(a) * wave()
		}
	}()
	return ch
}
