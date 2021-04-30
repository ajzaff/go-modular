package main

import (
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/control"
	"github.com/ajzaff/go-modular/drivers/otodriver"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
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

	var vca1 vca.VCA
	var vca2 vca.VCA

	cancel1 := mod.Patch(vca1.A(),
		osc.Pulse(.1, .1, osc.RangeLo, 0, .5),
		modular.NopProcessor(control.Voltage(0)))

	cancel2 := mod.Patch(mod.Send(0),
		&vca1,
		osc.Saw(.1, osc.Range32, osc.Fine(midi.StdTuning)),
		modular.NopProcessor(control.Voltage(69)))

	cancel3 := mod.Patch(vca2.A(),
		osc.Pulse(.1, .1, osc.RangeLo, 0, .5),
		modular.NopProcessor(control.Voltage(0)))

	cancel4 := mod.Patch(mod.Send(1),
		&vca2,
		osc.Saw(.1, osc.Range32, osc.Fine(midi.StdTuning)),
		modular.NopProcessor(control.Voltage(69)))

	time.Sleep(10 * time.Second)

	for _, cancel := range []func() error{
		cancel1,
		cancel2,
		cancel3,
		cancel4,
	} {
		if err := cancel(); err != nil {
			panic(err)
		}
	}
}
