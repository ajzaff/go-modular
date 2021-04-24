package main

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

const freq = 440

func main() {
	ctx := modular.New(
		modular.WithSampleRate(
			context.Background(), 96000), otodriver.New())

	modular.Send(ctx, 0, func() <-chan modular.V {
		ch := make(chan modular.V, modular.BufferSize(ctx))
		go func() {
			gate := osc.Square(ctx, 1, osc.RangeLo, 0, control.Voltage(ctx, 0))
			for v := range osc.Saw(ctx,
				1, osc.Range32, osc.Fine(midi.StdTuning),
				control.Voltage(ctx, float64(midi.Key(midi.StdTuning, freq)))) {
				if <-gate != 0 {
					ch <- v
				} else {
					ch <- 0
				}
			}
		}()
		return ch
	}())

}
