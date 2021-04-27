package main

import (
	"github.com/ajzaff/go-modular"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	drv := otodriver.New()
	ctx := &modular.Context{
		Options: modular.Options{
			SampleRate:       44100,
			BufferSize:       16,
			DriverBufferSize: 4096,
			SampleSize:       4096,
		},
		Driver: drv,
	}
	drv.InitContext(ctx)

	note := make(chan modular.V, 1)
	go func() {
		note <- 69
	}()

	ctx.SendReader(0, osc.SineReader(ctx, 1, osc.Range8,
		osc.Fine(midi.StdTuning), note))
}
