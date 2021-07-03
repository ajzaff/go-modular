package filter

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type LowPass struct {
	blockSize int
	rate      int
}

func (f *LowPass) SetConfig(cfg *modular.Config) {
	f.blockSize = cfg.BufferSize
	f.rate = cfg.SampleRate
}

func (f *LowPass) computeFilter() []complex128 {
	b := make([]complex128, f.blockSize) // FIXME: avoid allocations.
	for i := range b {
		t := -(float64(f.blockSize)-1)/2 + float64(i)
		v := (2 * real(b[i]) / float64(f.rate)) * sincfn(2*real(b[i])*t/float64(f.rate))
		b[i] = complex(v, 0)
	}
	for i, v := range window.Blackman(f.blockSize) {
		b[i] *= complex(v, 0)
	}
	return b
}

func sincfn(x float64) float64 {
	return math.Sin(math.Pi*x) / (math.Pi * x)
}

func (f *LowPass) Process(b []float32) {
	x := make([]complex128, len(b)) // FIXME: avoid allocations.
	for i, v := range b {
		x[i] = complex(float64(v), 0)
	}
	h := f.computeFilter()
	y := fft.Convolve(x, h)

	copy(x, y)
	for i, v := range x {
		b[i] = float32(real(v))
	}
}
