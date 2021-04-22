package main

import (
	"sync"
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	osc "github.com/ajzaff/go-modular/components/oscillator"
)

func main() {
	ctx := modular.NewContext(otodriver.New())

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		quit := make(chan struct{})
		go ctx.Send(0, osc.Sine(ctx, control.Voltage(ctx, .3), control.Voltage(ctx, 0), quit))
		time.Sleep(5 * time.Second)
		quit <- struct{}{}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Second)
		quit := make(chan struct{})
		go ctx.Send(1, osc.Sine(ctx, control.Voltage(ctx, .3), control.Voltage(ctx, 219.3), quit))
		time.Sleep(4 * time.Second)
		quit <- struct{}{}
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		quit := make(chan struct{})
		go ctx.Send(0, osc.Sine(ctx, control.Voltage(ctx, .3), control.Voltage(ctx, 344), quit))
		time.Sleep(3 * time.Second)
		quit <- struct{}{}
		wg.Done()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		quit := make(chan struct{})
		go ctx.Send(0, osc.Sine(ctx, control.Voltage(ctx, .3), control.Voltage(ctx, 547.77), quit))
		time.Sleep(2 * time.Second)
		quit <- struct{}{}
		wg.Done()
	}()

	wg.Wait()
}
