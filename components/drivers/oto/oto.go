package otodriver

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/ajzaff/go-modular"
	"github.com/hajimehoshi/oto"
)

type driver struct {
	*oto.Context
}

// New initializes a new Oto driver.
//
// Init should only be called once ever, so
// the driver should probably be singleton.
func New() modular.Driver {
	return &driver{}
}

const minBuffer = 4096

// Init initializes a new Oto driver.
//
// Init should only be called once.
// Init enforces a minimum buffer size of 1024 to ensure performance.
func (d *driver) Init(ctx context.Context) {
	sampleRate := modular.SampleRate(ctx)
	bufferSize := modular.DriverBufferSize(ctx)
	if bufferSize < minBuffer {
		bufferSize = minBuffer
	}
	oto, err := oto.NewContext(sampleRate, 2, 2, bufferSize)
	if err != nil {
		panic(err)
	}
	d.Context = oto
}

// Send outputs to the speaker using the Oto driver.
func (d *driver) Send(ch int, in <-chan modular.V) (n int64, err error) {
	switch player := d.NewPlayer(); ch {
	case 0:
		for v := range in {
			binary.Write(player, binary.LittleEndian, convert(float64(v)))
			binary.Write(player, binary.LittleEndian, int16(0))
			n++
		}
	case 1:
		for v := range in {
			binary.Write(player, binary.LittleEndian, int16(0))
			binary.Write(player, binary.LittleEndian, convert(float64(v)))
			n++
		}
	default:
		return 0, errors.New("otodriver.Send: only 2 stereo channels are supported [0, 1]")
	}
	return n, nil
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
