package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := make([]float32, 5*44100)

	noise := osc.NoiseOsc{}
	noise.SetConfig(cfg)
	noise.Process(b)

	for i, v := range b {
		b[i] = .1 * v
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.SendStereo().Process(b)
}
