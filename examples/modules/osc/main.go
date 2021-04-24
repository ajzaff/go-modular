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
		context.Background(), 96000), otodriver.New())
	modular.Send(ctx, 0, osc.Pulse(ctx, 1, osc.Range8,
		osc.Fine(midi.StdTuning),
		control.Voltage(ctx, .01),
		control.Voltage(ctx, float64(midi.Note(midi.A, 4)))))
}
