package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modules/filter"
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

	cfg.BufferSize = blockSize
	f := filter.LowPass{}
	f.SetConfig(cfg)
	f.UpdateFilter(func() float32 { return 440 })
	for i := 0; i+blockSize <= len(b); i += blockSize {
		f.Process(b[i : i+blockSize])
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
