package main

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
)

func main() {
	ctx := modular.New(context.Background(), otodriver.New())

	for v := range control.Voltage(ctx, +5) {
		println(v)
	}
}
