package filter

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type LowPass struct {
	Cutoff func(i int) float32

	blockSize int
	rate      int

	buffer []complex128
	filter []complex128
}

func (f *LowPass) SetConfig(cfg *modular.Config) {
	f.blockSize = cfg.BufferSize
	f.rate = cfg.SampleRate
}

func sincfn(x float32) float32 {
	return float32(math.Sin(math.Pi*float64(x)) / (math.Pi * float64(x)))
}

func (f *LowPass) computeFilter() {
	if f.filter == nil {
		f.filter = make([]complex128, f.blockSize)
	}
	for i := range f.filter {
		c := f.Cutoff(i)
		t := -(float32(f.blockSize)-1)/2 + float32(i)
		v := (2 * c / float32(f.rate)) * sincfn(2*c*t/float32(f.rate))
		f.filter[i] = complex(float64(v), 0)
	}
	for i, v := range window.Blackman(f.blockSize) {
		f.filter[i] *= complex(v, 0)
	}
}

func (f *LowPass) Process(b []float32) {
	if len(b) != f.blockSize {
		panic("filter.LowPass: Process called with wrong size block")
	}
	if f.buffer == nil {
		f.buffer = make([]complex128, len(b))
	}
	for i, v := range b {
		f.buffer[i] = complex(float64(v), 0)
	}
	f.computeFilter()
	y := fft.Convolve(f.buffer, f.filter)

	copy(f.buffer, y)
	for i, v := range f.buffer {
		b[i] = float32(real(v))
	}
}
