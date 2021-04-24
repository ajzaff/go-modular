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
	ctx := modular.New(context.Background(), otodriver.New())

	sine := osc.Sine(ctx, 1, osc.Range8, 0, control.Voltage(ctx, 69))
	mult := util.Mult(ctx, 2, sine)

	go modular.Send(ctx, 0, mult[0])
	modular.Send(ctx, 1, mult[1])
}
