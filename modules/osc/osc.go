// Package osc provides standard VCO and LFO waveforms.
package osc

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
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
	buf        []modular.V
	p          int
}

func (o *Osc) SetConfig(cfg *modular.Config) error {
	o.sampleRate = float64(cfg.SampleRate)
	return nil
}

func (o *Osc) BlockSize() int { return 512 }

func (o *Osc) Write(vs []modular.V) (n int, err error) {
	if o.buf == nil {
		o.buf = make([]modular.V, 512)
	}
	n = copy(o.buf[o.p:], vs)
	o.p += n
	return n, nil
}

func (o *Osc) Read(vs []modular.V) (n int, err error) {
	n = copy(vs, o.buf[:o.p])
	o.p -= n
	copy(o.buf, o.buf[n:])
	o.buf = o.buf[0:o.p:cap(o.buf)]
	for i, v := range vs[:n] {
		vs[i] = o.fn(v, float64(o.t), o.sampleRate)
		o.t++
	}
	return n, nil
}

func (o *Osc) ReadStream() modular.V {
	var v modular.V
	if o.p < len(o.buf) {
		v = o.buf[o.p]
		o.p++
	}
	return o.fn(v, float64(o.t), o.sampleRate)
}

// Sine outputs an sine audio wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Sine(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) modular.V {
			return modular.V(a) * modular.V(math.Sin(2*math.Pi*float64(t)*Tone(r, fine+float64(x))/float64(sampleRate)))
		},
	}
}

// Triangle outputs an triangle wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Triangle(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			return modular.V(a) * modular.V(2/math.Pi*math.Asin(math.Sin(2*math.Pi*float64(t)*Tone(r, fine+float64(x))/float64(sampleRate))))
		},
	}
}

// Saw outputs an sawtooth wave from the linear signal and parameters.
//
// Linear signal lin conforms to the real midi scale (one volt per octave).
func Saw(a Polarity, r Range, fine float64) *Osc {
	return &Osc{
		fn: func(x modular.V, t, sampleRate float64) (v modular.V) {
			return modular.V(a) * modular.V(2/math.Pi*math.Atan(math.Tan(math.Pi*float64(t)*Tone(r, fine+float64(x))/float64(sampleRate))))
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
			if math.Mod(float64(t)*Tone(r, fine+float64(x))/float64(sampleRate), 2) < 2*float64(w) {
				v = modular.V(a) + modular.V(c)
			} else {
				v = modular.V(-a) + modular.V(c)
			}
			return
		},
	}
}
