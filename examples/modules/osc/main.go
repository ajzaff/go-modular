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
	w := osc.Sine(.1, osc.Range8, osc.Fine(midi.StdTuning))
	w.Voltage = func() float32 {
		return 69. / 12
	}
	w.SetConfig(cfg)
	w.Process(b)

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
