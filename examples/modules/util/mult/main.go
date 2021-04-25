package main

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
	"github.com/ajzaff/go-modular/modules/util"
)

func main() {
	ctx := modular.New(
		modular.WithBufferSize(context.Background(), 10000),
		otodriver.New())

	sine := osc.Sine(ctx, 1, osc.Range8, 0, control.Voltage(ctx, 69))
	mult := util.Mult(ctx, 2, sine)

	go func() {
		_, err := modular.Send(ctx, 0, mult[0])
		panic(err)
	}()
	_, err := modular.Send(ctx, 1, mult[1])
	panic(err)
}
