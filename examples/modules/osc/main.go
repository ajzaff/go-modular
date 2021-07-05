package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := make([]float32, 5*44100)
	for i := range b {
		b[i] = 69. / 12
	}

	wave := osc.Sine(.1, osc.Range8, osc.Fine(midi.StdTuning))
	wave.SetConfig(cfg)
	wave.Process(b)

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
