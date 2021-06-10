package main

import (
	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/midi"
	"github.com/ajzaff/go-modular/modules/osc"
	"github.com/ajzaff/go-modular/modules/output/otoplayer"
)

func main() {
	cfg := modular.New()
	buf := make([]float32, 44100*5)

	sine := osc.Sine(.1, osc.Range8, osc.Fine(midi.StdTuning))
	sine.SetConfig(cfg)
	sine.Process(modular.Block{
		Buf: buf,
	})

	// var wg sync.WaitGroup

	// for i := 0; i < 20; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		for j := range buf[i*len(buf)/20 : (i+1)*len(buf)/20] {
	// 			buf[j] = 69. / 12
	// 		}
	// 		sine.Process(buf[i*len(buf)/20 : (i+1)*len(buf)/20])
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()

	oto := otoplayer.New()
	oto.SetConfig(cfg)
	oto.Send(0).Process(modular.Block{
		Buf: buf,
	})
}
