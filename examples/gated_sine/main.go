package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := make([]float32, 10*44100)
	for i := range b {
		b[i] = 69. / 12
	}

	w := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	w.Voltage = func() float32 {
		return 69. / 12
	}
	w.SetConfig(cfg)
	w.Process(b)

	lfo := osc.Pulse(.5, .5, osc.Range64, 0, .5)
	lfo.SetConfig(cfg)

	for i, v := range b {
		b[i] = v * lfo.Next()
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
