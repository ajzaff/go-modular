// Package osc provides standard VCO and LFO waveforms.
package osc

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modio"
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
	RangeLo Range = iota // 1hz
	Range32              // 2hz
	Range16              // 4hz
	Range8               // 8hz
	Range4               // 16hz
	Range2               // 32hz
)

// Tone returns the tone frequency for the range and fine tuning.
func Tone(r Range, fine float64) float64 {
	return math.Pow(2, float64(r)+fine/12)
}

// Fine returns the fine tuning constant to tune the oscillators to t at Range8.
func Fine(t midi.Tuning) float64 {
	return 12*math.Log2(t.A4Hz()) - 105
}

type Osc struct {
	fn         func(v modular.V, t, sampleRate float64) modular.V
	t          int
	sampleRate float64
	buf        modio.Buffer
}

func (o *Osc) SetConfig(cfg *modular.Config) error {
	o.sampleRate = float64(cfg.SampleRate)
	return nil
}

func (o *Osc) BlockSize() int { return 0 }

func (o *Osc) Write(vs []modular.V) (n int, err error) {
	return o.buf.Write(vs)
}

func (o *Osc) Read(vs []modular.V) (n int, err error) {
	n, err = o.buf.Read(vs)
	for i, v := range vs[:n] {
		vs[i] = o.fn(v, float64(o.t), o.sampleRate)
		o.t++
	}
	return n, err
}

// Sine outputs an sine audio wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Sine(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			wavelen := float64(sampleRate) / Tone(r, fine+float64(x))
			v = modular.V(a) * modular.V(math.Sin(2*math.Pi*float64(t)/wavelen))
			return
		},
	}
}

// Triangle outputs an triangle wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Triangle(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			length := float64(sampleRate) / Tone(r, fine+float64(x))
			v = modular.V(a) * modular.V(2/math.Pi*math.Asin(math.Sin(2*math.Pi*float64(t)/length)))
			return
		},
	}
}

// Saw outputs an sawtooth wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Saw(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			length := float64(sampleRate) / Tone(r, fine+float64(x))
			v = modular.V(a) * modular.V(2/math.Pi*math.Atan(math.Tan(math.Pi*float64(t)/length)))
			return
		},
	}
}

// Square outputs an square wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Square(a Polarity, r Range, fine float64) *Osc {
	return Pulse(a, 0, r, fine, .5)
}

// Pulse outputs an pulse wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
// Pulse width w are in the range 0 to 1.
//
// TODO: Support PWM (again).
func Pulse(a Polarity, c float64, r Range, fine float64, w modular.V) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			length := float64(sampleRate) / Tone(r, fine+float64(x))
			if math.Mod(float64(t)/length, 2) < 2*float64(w) {
				v = modular.V(a) + modular.V(c)
			} else {
				v = modular.V(-a) + modular.V(c)
			}
			return
		},
	}
}
