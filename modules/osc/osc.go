// Package osc provides standard VCO and LFO waveforms.
package osc

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/mathmod"
)

// Polarity controls the polarity of waveform functions.
//
// Negative "inverts" the wave, while positive maintains
// the true shape. Conveniently, non-pole values can be
// used to set the amplitude.
type Polarity float32

const (
	Negative Polarity = 2*iota - 1 // inverted
	Positive                       // regular
)

// Range presents a pipe organ length setting.
// The zero value is LFO and higher values are
// octaves at 32hz doubling at each setting.
type Range int

const (
	Range64 Range = iota // 1hz
	Range32              // 2hz
	Range16              // 4hz
	Range8               // 8hz
	Range4               // 16hz
	Range2               // 32hz
)

// Tone returns the tone frequency for the range and fine tuning.
func Tone(r Range, fine float32) float32 {
	return float32(math.Pow(2, float64(r)+float64(fine)))
}

// Fine returns the fine tuning constant to tune the oscillators to t at Range8.
func Fine(t midi.Tuning) float32 {
	return 12*float32(math.Log2(float64(t.A4Hz()))) - 105
}

// Osc is a simple wave oscillator.
type Osc struct {
	mathmod.Func2

	sampleRate float32
	Phase      float32
}

func (o *Osc) SetConfig(cfg *modular.Config) error {
	o.sampleRate = float32(cfg.SampleRate)
	return nil
}

// Sine outputs an sine audio wave from the linear signal and parameters.
func Sine(a Polarity, r Range, fine float32) *Osc {
	osc := &Osc{}
	osc.Func2 = func(i int, x float32) float32 {
		return float32(a) * float32(math.Sin(2*math.Pi*(float64(osc.Phase)+float64(i))*float64(Tone(r, fine+x))/float64(osc.sampleRate)))
	}
	return osc
}

// Triangle outputs an triangle wave from the linear signal and parameters.
func Triangle(a Polarity, r Range, fine float32) *Osc {
	osc := &Osc{}
	osc.Func2 = func(i int, x float32) float32 {
		return float32(a) * float32(2/math.Pi*math.Asin(math.Sin(2*math.Pi*(float64(osc.Phase)+float64(i))*float64(Tone(r, fine+x))/float64(osc.sampleRate))))
	}
	return osc
}

// Saw outputs an sawtooth wave from the linear signal and parameters.
func Saw(a Polarity, r Range, fine float32) *Osc {
	osc := &Osc{}
	osc.Func2 = func(i int, x float32) float32 {
		return float32(a) * float32(2/math.Pi*math.Atan(math.Tan(math.Pi*(float64(osc.Phase)+float64(i))*float64(Tone(r, fine+x))/float64(osc.sampleRate))))
	}
	return osc
}

// Square outputs an square wave from the linear signal and parameters.
func Square(a Polarity, r Range, fine float32) *Osc {
	return Pulse(a, 0, r, fine, .5)
}

// Pulse outputs an pulse wave from the linear signal and parameters.
//
// Pulse width w is in the range 0 to 1.
func Pulse(a Polarity, c float32, r Range, fine float32, w float32) *Osc {
	osc := &Osc{}
	osc.Func2 = func(i int, x float32) float32 {
		if math.Mod((float64(osc.Phase)+float64(i))*float64(Tone(r, fine+x))/float64(osc.sampleRate), 2) < 2*float64(w) {
			return float32(a) + float32(c)
		}
		return float32(-a) + float32(c)
	}
	return osc
}
