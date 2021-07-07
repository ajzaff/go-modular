package filter

import (
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modio"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type LowPass struct {
	blockSize int
	rate      int

	filter []complex128
	fft    modio.FFT
}

func (f *LowPass) SetConfig(cfg *modular.Config) {
	f.blockSize = cfg.BufferSize
	f.rate = cfg.SampleRate
}

func sincfn(x float32) float32 {
	return float32(math.Sin(math.Pi*float64(x)) / (math.Pi * float64(x)))
}

// Adapted from https://ccrma.stanford.edu/~jos/sasp/Example_1_Low_Pass_Filtering.html.
func (f *LowPass) UpdateFilter(cutoff func() float32) {
	filter := f.filter
	if f.filter == nil {
		filter = make([]complex128, f.blockSize)
	}
	for i := range filter {
		c := cutoff()
		t := -(float32(f.blockSize)-1)/2 + float32(i)
		v := (2 * c / float32(f.rate)) * sincfn(2*c*t/float32(f.rate))
		filter[i] = complex(float64(v), 0)
	}
	for i, v := range window.Blackman(f.blockSize) {
		filter[i] *= complex(v, 0)
	}
	f.filter = fft.FFT(filter)
}

func (f *LowPass) Process(b []float32) {
	if len(b) != f.blockSize {
		panic("filter.LowPass: Process called with wrong size block")
	}
	if f.filter == nil {
		return
	}
	f.fft.Reset()
	f.fft.StoreFFT(b)
	f.fft.UpdateAll(func(i int, v complex128) complex128 {
		return v * f.filter[i]
	})
	f.fft.Process(b)
}
