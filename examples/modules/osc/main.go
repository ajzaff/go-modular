package main

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	ctx := modular.New(modular.WithSampleRate(
		context.Background(), 44000), otodriver.New())
	modular.Send(ctx, 0, osc.Sine(ctx, 1, osc.Range8,
		control.V(osc.Fine(midi.StdTuning)),
		control.Voltage(ctx, control.V(midi.Note(midi.A, 4)))))
}
