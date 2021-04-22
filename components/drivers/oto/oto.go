package otodriver

import (
	"encoding/binary"
	"errors"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/utility"
	"github.com/ajzaff/go-modular/sampleio"
	"github.com/hajimehoshi/oto"
)

type driver struct {
	ctx *modular.Context
	*oto.Context
}

// New initializes a new Oto driver.
//
// New should only be called once.
func New() modular.Driver {
	return &driver{}
}

// Init initializes a new Oto driver.
//
// Init should only be called once.
func (d *driver) Init(ctx *modular.Context) {
	oto, err := oto.NewContext(ctx.SampleRate, 2, 2, ctx.BufferSize)
	if err != nil {
		panic(err)
	}
	d.ctx = ctx
	d.Context = oto
}

// Send outputs to the speaker using the Oto driver.
func (d *driver) Send(ch int, in <-chan modular.V) (n int64, err error) {
	return sampleio.Copy(d.ctx, &otoWriter{ch, d.Context.NewPlayer()}, utility.Reader(in))
}

type otoWriter struct {
	ch     int
	player *oto.Player
}

func (w *otoWriter) Write(vs []modular.V) (n int, err error) {
	switch w.ch {
	case 0:
		for _, v := range vs {
			binary.Write(w.player, binary.LittleEndian, convert(float64(v)))
			binary.Write(w.player, binary.LittleEndian, int16(0))
		}
		return len(vs), nil
	case 1:
		for _, v := range vs {
			binary.Write(w.player, binary.LittleEndian, int16(0))
			binary.Write(w.player, binary.LittleEndian, convert(float64(v)))
		}
		return len(vs), nil
	default:
		return 0, errors.New("otodriver.Write: only stereo channels are supported [0, 1]")
	}
}

func convert(x float64) int16 {
	if x < -1 {
		x = -1
	}
	if x > 1 {
		x = 1
	}
	return int16((-1 << 15) + (x+1)/2*(1<<16-1))
}
