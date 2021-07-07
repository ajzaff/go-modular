package main

import (
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/adsr"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := make([]float32, 5*44100)
	w := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	w.SetConfig(cfg)
	w.Process(b)

	g := adsr.New(time.Second, time.Second, .5, time.Second)
	g.SetConfig(cfg)
	g.SetSustain(time.Second)
	g.Process(b)

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.PlayStereo(b)
}
