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
	for i := range b {
		b[i] = 69. / 12
	}

	noise := osc.Noise(.1)
	noise.SetConfig(cfg)
	noise.Process(b)

	// w := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	// w.SetConfig(cfg)
	// w.Process(b)

	fft := modio.FFT{}

	for i := 0; i+blockSize <= len(b); i += blockSize {
		block := b[i : i+blockSize]

		fft.Reset()
		fft.StoreFFT(block)
		fft.UpdateAll(func(i int, v complex128) complex128 {
			if i >= 100 && i <= 500 {
				return v
			}
			return 0
		})
		fft.Process(block)
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
