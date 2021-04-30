package main

import (
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/control"
	"github.com/ajzaff/go-modular/drivers/otodriver"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
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

	cancel := mod.Patch(mod.Send(0),
		osc.Sine(.1, osc.Range8, osc.Fine(midi.StdTuning)),
		modular.NopProcessor(control.Voltage(69)),
	)

	time.Sleep(5 * time.Second)

	if err := cancel(); err != nil {
		panic(err)
	}
}
