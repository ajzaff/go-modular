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
		b[i] = w.Func2(i, 45./12)/3 + w.Func2(i, 64./12)/3 + w.Func2(i, 73./12)/3
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.SendStereo().Process(b)
}
