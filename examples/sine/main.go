package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()
	buf := make([]float32, 44100*5)

	sine := osc.Sine(.1, osc.Range8, osc.Fine(midi.StdTuning))
	sine.SetConfig(cfg)
	sine.Process(buf)

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.Send(0).Process(buf)
}
