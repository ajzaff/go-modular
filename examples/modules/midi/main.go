package main

import (
	"context"
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	"github.com/ajzaff/go-modular/modules/adsr"
	midimodule "github.com/ajzaff/go-modular/modules/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
	"github.com/ajzaff/go-modular/modules/vca"
)

func main() {
	ctx := modular.New(
		modular.WithDriverBufferSize(context.Background(), 4096),
		otodriver.New())

	gate, key, vel, err := midimodule.Interface(ctx, 1, 0)
	if err != nil {
		panic(err)
	}

	go func() {
		w := osc.Saw(ctx, 1,
			osc.Range8, osc.Fine(midi.StdTuning),
			control.Latch(ctx, key))

		eg := adsr.Envelope(ctx,
			/* a */ 0,
			/* d */ 100*time.Millisecond,
			/* s */ 0.5,
			/* r */ 10*time.Millisecond,
			control.Latch(ctx, gate))

		amp := vca.VCA(ctx, eg, w)

		modular.Send(ctx, 0, amp)
	}()
	for range vel {
	}
}
