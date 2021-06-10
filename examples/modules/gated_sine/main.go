package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()

	b := modular.Block{
		Buf: make([]float32, 10*44100),
	}

	for i := range b.Buf {
		b.Buf[i] = 69. / 12
	}

	lb := modular.Block{
		Buf: make([]float32, 10*44100),
	}
	lfo := osc.Pulse(.1, .1, osc.Range64, 0, .5)
	lfo.SetConfig(cfg)
	lfo.Process(lb)

	wave := osc.Sine(.1, osc.Range16, osc.Fine(midi.StdTuning))
	wave.SetConfig(cfg)
	wave.Process(b)

	for i := range lb.Buf {
		b.Buf[i] = lb.Buf[i] * b.Buf[i]
	}

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.Send(0).Process(b)
}
