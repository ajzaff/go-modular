package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	ctx := modular.WithSampleRate(
		modular.NewContext(otodriver.New()), 44000)
	modular.Send(ctx, 0, osc.Sine(ctx, 1, osc.Range8,
		control.V(osc.Fine(midi.StdTuning)),
		control.Voltage(ctx, control.V(midi.Note(midi.A, 4)))))
}
