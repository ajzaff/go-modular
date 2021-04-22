package main

import (
	"github.com/ajzaff/go-modular"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	ctx := modular.NewContext(otodriver.New())
	modular.Send(ctx, 0, osc.Saw(ctx, osc.RangeLow, 5))
}
