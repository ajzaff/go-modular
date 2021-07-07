package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modio"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

const blockSize = 1024

func main() {
	cfg := modular.New()

	b := make([]float32, 5*44100)
	noise := osc.Noise(.1)
	noise.SetConfig(cfg)
	noise.Process(b)

	// w := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	// w.SetConfig(cfg)
	// w.Process(b)

	fft := modio.FFT{}

	for i := 0; i+blockSize <= len(b); i += blockSize {
		block := b[i : i+blockSize]

		h := make([]float32, blockSize)
		for i := blockSize / 2; i < blockSize; i++ {
			h[i] = 1
		}
		hc := fft.Compute(h)

		fft.Reset()
		fft.StoreFFT(block)
		fft.UpdateAll(func(i int, v complex128) complex128 {
			return v * hc[i]
		})
		fft.Process(block)
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
