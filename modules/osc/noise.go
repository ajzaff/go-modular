package osc

import (
	"math"

	"github.com/ajzaff/go-modular"
)

const noiseSeed = 1260667865

// Xorshift32 from p. 4 of Marsaglia, "Xorshift RNGs"
type NoiseOsc struct {
	State uint32
	A     Polarity
	C     float32
}

func Noise(a Polarity) *NoiseOsc {
	return &NoiseOsc{A: a}
}

func (*NoiseOsc) SetConfig(*modular.Config) {}

func (o *NoiseOsc) Next() float32 {
	v, x := nextRand(o.State)
	o.State = x
	return v*float32(o.A) + o.C
}

func nextRand(x uint32) (v float32, state uint32) {
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	return 2*(float32(x)/math.MaxUint32) - 1, x
}

func (o *NoiseOsc) Process(b []float32) {
	x := o.State
	if x == 0 {
		x = noiseSeed
	}
	for i := range b {
		var v float32
		v, x = nextRand(x)
		b[i] = v*float32(o.A) + o.C
	}
	o.State = x
}
