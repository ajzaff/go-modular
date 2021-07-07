package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modules/filter"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

const (
	blockSize  = 1024
	sampleRate = 44100
)

func main() {
	cfg := modular.New()
	cfg.SampleRate = sampleRate
	cfg.BufferSize = blockSize

	b := make([]float32, 7*sampleRate)

	noise := osc.Noise(.1)
	noise.SetConfig(cfg)
	noise.Process(b)

	f := filter.LowPass{}
	f.SetConfig(cfg)

	var i float32
	f.SetCutoff(func() float32 {
		defer func() { i++ }()
		x := sampleRate / 2 * i / (5 * sampleRate)
		if x > sampleRate/2 {
			return 0
		}
		return sampleRate/2 - x
	})

	for i := 0; i+blockSize <= len(b); i += blockSize {
		f.Process(b[i : i+blockSize])
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
