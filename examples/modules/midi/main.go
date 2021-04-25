package main

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	midimodule "github.com/ajzaff/go-modular/modules/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	ctx := modular.New(context.Background(), otodriver.New())

	gate, key, vel, err := midimodule.Interface(ctx, 1, 0)
	if err != nil {
		panic(err)
	}

	go func() {
		for range gate {
		}
	}()
	go func() {
		w := osc.Saw(ctx, .5,
			osc.Range8, osc.Fine(midi.StdTuning),
			control.Latch(ctx, key))
		// mult := util.Mult(ctx, 2, w)
		// go modular.Send(ctx, 0, mult[0])
		// modular.Send(ctx, 1, mult[1])
		modular.Send(ctx, 0, w)
	}()
	for range vel {
	}
}
