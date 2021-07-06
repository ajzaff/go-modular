package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	w := osc.Sine(.5, osc.Range16, osc.Fine(midi.StdTuning))
	w.SetConfig(cfg)

	b := make([]float32, 5*44100)
	for i := range b {
		b[i] = w.Func(45./12)/3 + w.Func(64./12)/3 + w.Func(73./12)/3
		w.Advance(1)
	}

	oto := otoplayer.New()
	defer oto.Close()

	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
