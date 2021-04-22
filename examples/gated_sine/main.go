package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	osc "github.com/ajzaff/go-modular/components/oscillator"
)

const freq = 440

func main() {
	ctx := modular.NewContext(otodriver.New())

	i := 0
	modular.Send(ctx, 0, osc.Sine(ctx, control.Func(ctx, func() control.V {
		if i++; (i/modular.SampleRate(ctx))%2 == 0 {
			return 1
		}
		return 0
	}), control.Voltage(ctx, freq) /* quit = */, make(chan struct{})))
}
