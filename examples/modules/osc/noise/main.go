package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := make([]float32, 5*44100)

	noise := osc.Noise(.1)
	noise.SetConfig(cfg)
	noise.Process(b)

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
