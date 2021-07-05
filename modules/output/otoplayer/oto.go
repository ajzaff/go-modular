package otoplayer

import (
	"encoding/binary"
	"sync"

	"github.com/ajzaff/go-modular"
	"github.com/hajimehoshi/oto"
)

type Context struct {
	ctx *oto.Context
	mu  sync.Mutex
}

// New creates a new Oto output context.
func New() *Context {
	return &Context{}
}

func (d *Context) SetConfig(cfg *modular.Config) error {
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

func (d *Context) NewStereoPlayer() *StereoPlayer {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ctx == nil {
		panic("otoplayer.SendStereo called before SetConfig")
	}
	return &StereoPlayer{d.ctx.NewPlayer()}
}

func (d *Context) PlayStereo(b []float32) {
	p := d.NewStereoPlayer()
	defer p.Close()
	p.Process(b)
}

func (d *Context) NewPlayer(ch int) *Player {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ctx == nil {
		panic("otoplayer.SendStereo called before SetConfig")
	}
	return &Player{d.ctx.NewPlayer(), ch}
}

func (d *Context) Play(ch int, b []float32) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ctx == nil {
		panic("otoplayer.Send called before SetConfig")
	}
	p := &Player{d.ctx.NewPlayer(), ch}
	defer p.Close()
	p.Process(b)
}

func (d *Context) Close() error {
	if d.ctx == nil {
		return nil
	}
	return d.ctx.Close()
}

type Player struct {
	player *oto.Player
	ch     int
}

// Send outputs to the speaker using the Oto driver.
func (d *Player) Process(b []float32) {
	if d.ch == 1 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	for _, v := range b[:len(b)-1] {
		binary.Write(d.player, binary.LittleEndian, convert(v))
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
	binary.Write(d.player, binary.LittleEndian, convert(b[len(b)-1]))
	if d.ch == 0 {
		binary.Write(d.player, binary.LittleEndian, uint16(0))
	}
}

func (d *Player) Close() error {
	return d.player.Close()
}

type StereoPlayer struct {
	player *oto.Player
}

// Send outputs to the speaker using the Oto driver.
func (d *StereoPlayer) Process(b []float32) {
	for _, v := range b {
		binary.Write(d.player, binary.LittleEndian, convert(v))
		binary.Write(d.player, binary.LittleEndian, convert(v))
	}
}

func (d *StereoPlayer) Close() error {
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
