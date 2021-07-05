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

	wave := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	wave.SetConfig(cfg)
	wave.Process(b)

	lfo := osc.Pulse(.1, .1, osc.Range64, 0, .5)
	lfo.SetConfig(cfg)

	for i := range b {
		b[i] *= lfo.Func2(i, 0)
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
