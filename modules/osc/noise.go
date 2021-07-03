package osc

import (
	"math"

	"github.com/ajzaff/go-modular"
)

const noiseSeed = 1260667865

type NoiseOsc struct {
	State uint32
}

func Noise() *NoiseOsc {
	return &NoiseOsc{}
}

func (*NoiseOsc) SetConfig(*modular.Config) {}

func (o *NoiseOsc) Process(b []float32) {
	x := o.State
	if x == 0 {
		x = noiseSeed
	}
	for i := range b {
		// Xorshift32 from p. 4 of Marsaglia, "Xorshift RNGs"
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		b[i] = 2*(float32(x)/math.MaxUint32) - 1
	}
	o.State = x
}
