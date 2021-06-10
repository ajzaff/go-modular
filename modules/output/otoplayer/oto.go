package otoplayer

import (
	"encoding/binary"
	"sync"

	"github.com/ajzaff/go-modular"
	"github.com/hajimehoshi/oto"
)

type Player struct {
	ctx *oto.Context
	mu  sync.Mutex
}

// New creates a new Oto output.
//
// New may panic if called again before the driver is Closed.
func New() *Player {
	return &Player{}
}

func (d *Player) SetConfig(cfg *modular.Config) error {
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

func (d *Player) Send(ch int) modular.Processor {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ctx == nil {
		panic("otoplayer.Send called before SetConfig")
	}
	return &playerProcessor{d.ctx.NewPlayer(), ch}
}

func (d *Player) Close() error {
	if d.ctx == nil {
		return nil
	}
	return d.ctx.Close()
}

type playerProcessor struct {
	player *oto.Player
	ch     int
}

// Send outputs to the speaker using the Oto driver.
func (d *playerProcessor) Process(b modular.Block) {
	if d.ch == 1 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	for _, v := range b.Buf[:len(b.Buf)-1] {
		binary.Write(d.player, binary.LittleEndian, convert(v))
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	binary.Write(d.player, binary.LittleEndian, convert(b.Buf[len(b.Buf)-1]))
	if d.ch == 0 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
}

func (d *playerProcessor) Close() error {
	return d.player.Close()
}

func convert(x float32) int16 {
	if x < -1 {
		x = -1
	}
	if x > 1 {
		x = 1
	}
	return int16((-1 << 15) + (x+1)/2*(1<<16-1))
}
