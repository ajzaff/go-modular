package otodriver

import (
	"encoding/binary"
	"sync"

	"github.com/ajzaff/go-modular"
	"github.com/hajimehoshi/oto"
)

type Driver struct {
	ctx *oto.Context
	mu  sync.Mutex
}

// New creates a new Oto driver.
//
// New may panic if called again before the driver is Closed.
func New() *Driver {
	return &Driver{}
}

func (d *Driver) SetConfig(cfg *modular.Config) error {
	if err := d.Close(); err != nil {
		return err
	}
	oto, err := oto.NewContext(cfg.SampleRate, 2, 2, cfg.DriverBufferSize)
	if err != nil {
		panic(err)
	}
	d.ctx = oto
	return nil
}

func (d *Driver) Send(ch int) modular.WriteCloser {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ctx == nil {
		if err := d.SetConfig(modular.NewConfig()); err != nil {
			panic(err)
		}
	}
	return &driverWriter{d.ctx.NewPlayer(), ch}
}

func (d *Driver) Close() error {
	if d.ctx == nil {
		return nil
	}
	return d.ctx.Close()
}

type driverWriter struct {
	player *oto.Player
	ch     int
}

// Send outputs to the speaker using the Oto driver.
func (d *driverWriter) Write(vs []modular.V) (n int, err error) {
	if d.ch == 1 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	for _, v := range vs[:len(vs)-1] {
		binary.Write(d.player, binary.LittleEndian, convert(float64(v)))
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	binary.Write(d.player, binary.LittleEndian, convert(float64(vs[len(vs)-1])))
	if d.ch == 0 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	return len(vs), nil
}

func (d *driverWriter) Close() error {
	return d.player.Close()
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
