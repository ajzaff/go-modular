package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/drivers/otodriver"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/util"
)

func main() {
	mod, err := modular.New(otodriver.New())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			panic(err)
		}
	}()

	var mult util.Mult

	w := osc.Saw(.1, osc.Range8, 0)

	mod.Patch(&mult, modular.NopProcessor(w))
	mod.Patch(mod.Send(0), modular.NopProcessor(mult.New()))
	mod.Patch(mod.Send(1), modular.NopProcessor(mult.New()))
}
