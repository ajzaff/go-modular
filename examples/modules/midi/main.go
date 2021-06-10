package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	midimodule "github.com/ajzaff/go-modular/modules/midi"
	"github.com/ajzaff/go-modular/modules/osc"
)

func main() {
	cfg := modular.New()

	mid, err := midimodule.New(1, 0)
	if err != nil {
		panic(err)
	}
	mid.SetConfig(cfg)

	wave := osc.Saw(.5, osc.Range8, osc.Fine(midi.StdTuning))
	wave.SetConfig(cfg)

	// FIXME:

}
