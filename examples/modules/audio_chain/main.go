package main

import (
	"context"
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	"github.com/ajzaff/go-modular/modules/adsr"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
	"github.com/ajzaff/go-modular/modules/vca"
)

func main() {
	ctx := modular.New(context.Background(), otodriver.New())

	w := osc.Sine(ctx, 1, osc.Range8,
		osc.Fine(midi.StdTuning),
		control.Voltage(ctx, float64(midi.Note(midi.A, 4))))

	gate := osc.Square(ctx,
		1, osc.RangeLo, 0,
		control.Voltage(ctx, -12))

	eg := adsr.Envelope(ctx,
		/* a */ 10*time.Millisecond,
		/* d */ 100*time.Millisecond,
		/* s */ 0.1,
		/* r */ 10*time.Millisecond, gate)

	amp := vca.VCA(ctx, eg, w)

	modular.Send(ctx, 0, amp)
}
