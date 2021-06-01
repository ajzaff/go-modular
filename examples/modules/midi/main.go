package main

import (
	"sync"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/drivers/otodriver"
	"github.com/ajzaff/go-modular/midi"
	midimodule "github.com/ajzaff/go-modular/modules/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/util"
	"github.com/ajzaff/go-modular/modules/vca"
)

func main() {
	mod, err := modular.New(otodriver.New())
	if err != nil {
		panic(err)
	}
	defer mod.Close()

	var mult util.Mult
	var amp vca.VCA

	mid, err := midimodule.New(1, 0)
	if err != nil {
		panic(err)
	}

	_ = mult

	gate, key := mid.GateKey()

	patches := map[string]func() error{
		"amp_wave": mod.Patch(&amp,
			osc.Saw(.5, osc.Range8, osc.Fine(midi.StdTuning)),
			modular.NopProcessor(key)),
		"amp_gate": mod.Patch(amp.A(), modular.NopProcessor(gate)),
		"send_amp": mod.Patch(mod.Send(0), &amp),
	}

	defer func() {
		for _, cancel := range patches {
			cancel()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
