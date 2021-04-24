package main

import (
	"sync"
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	otodriver "github.com/ajzaff/go-modular/components/drivers/oto"
	"github.com/ajzaff/go-modular/components/midi"
	osc "github.com/ajzaff/go-modular/modules/oscillator"
)

func main() {
	ctx := modular.NewContext(otodriver.New())

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		go modular.Send(ctx, 0, osc.Sine(ctx, .2,
			osc.Range8, control.V(osc.Fine(midi.StdTuning)),
			control.Voltage(ctx, control.V(midi.Note(midi.A, 4)))))
		time.Sleep(5 * time.Second)
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Second)
		go modular.Send(ctx, 1, osc.Sine(ctx, .2,
			osc.Range8, control.V(osc.Fine(midi.StdTuning)),
			control.Voltage(ctx, control.V(midi.Note(midi.C, 5)))))
		time.Sleep(4 * time.Second)
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		go modular.Send(ctx, 0, osc.Sine(ctx, .2,
			osc.Range8, control.V(osc.Fine(midi.StdTuning)),
			control.Voltage(ctx, control.V(midi.Note(midi.E, 5)))))
		time.Sleep(3 * time.Second)
		wg.Done()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		go modular.Send(ctx, 1, osc.Sine(ctx, .2,
			osc.Range8, control.V(osc.Fine(midi.StdTuning)),
			control.Voltage(ctx, control.V(midi.Note(midi.G, 5)))))
		time.Sleep(2 * time.Second)
		wg.Done()
	}()

	wg.Wait()
}
