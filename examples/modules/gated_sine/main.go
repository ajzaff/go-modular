package main

import (
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/control"
	"github.com/ajzaff/go-modular/drivers/otodriver"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/util"
	"github.com/ajzaff/go-modular/modules/vca"
)

func main() {
	mod, err := modular.New(otodriver.New())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			panic(err)
		}
	}()

	var mult util.Mult
	left, right := mult.New(), mult.New()

	var amp [2]vca.VCA
	patches := map[string]func() error{
		"lfo_mult": mod.Patch(&mult,
			osc.Pulse(.1, .1, osc.RangeLo, 0, .5),
			modular.NopProcessor(control.Voltage(0))),
		"left_lfo": mod.Patch(amp[0].A(), modular.NopProcessor(left)),
		"left_audio": mod.Patch(mod.Send(0),
			&amp[0],
			osc.Saw(.1, osc.Range32, osc.Fine(midi.StdTuning)),
			modular.NopProcessor(control.Voltage(69))),
		"right_lfo": mod.Patch(amp[1].A(), modular.NopProcessor(right)),
		"right_audio": mod.Patch(mod.Send(1),
			&amp[1],
			osc.Saw(.1, osc.Range32, osc.Fine(midi.StdTuning)),
			modular.NopProcessor(control.Voltage(69))),
	}

	time.Sleep(10 * time.Second)

	for _, cancel := range patches {
		if err := cancel(); err != nil {
			panic(err)
		}
	}
}
