package midi

import (
	"fmt"

	"github.com/ajzaff/go-modular"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

// Interface returns midi CVs for the stream of midi messages
// on the single midi input channel (i, ch).
//
// Interface is unbuffered to minimize trigger latency.
type Interface struct {
	gate modular.V
	in   midi.In
	ch   uint8
	drv  *rtmididrv.Driver
}

// New creates a new midi interface on input i MIDI channel ch.
func New(i, ch uint8) (iface *Interface, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("midi.New: %v", err)
		}
	}()

	drv, err := rtmididrv.New()
	if err != nil {
		return nil, err
	}

	ins, err := drv.Ins()
	if len(ins) <= int(i) {
		return nil, err
	}

	in := ins[i]
	if err := in.Open(); err != nil {
		return nil, err
	}

	iface = &Interface{ch: ch, in: in}
	return iface, nil
}

func (i *Interface) UpdateConfig(cfg *modular.Config) error {
	return nil
}

func (i *Interface) GateKey() (gate, key modular.Reader) {
	ch := i.ch
	g, k := &midiGateReader{}, &midiKeyReader{}
	rd := reader.New(reader.NoLogger())
	rd.Channel.NoteOn = func(p *reader.Position, channel uint8, key uint8, velocity uint8) {
		if channel != ch {
			return
		}
		k.key = modular.V(key)
		g.gate = 1
	}
	rd.Channel.NoteOff = func(p *reader.Position, channel uint8, key uint8, velocity uint8) {
		if channel != ch {
			return
		}
		g.gate = 0
	}
	if err := rd.ListenTo(i.in); err != nil {
		panic(fmt.Errorf("midi.Interface.Key: %v", err))
	}
	return g, k
}

type midiKeyReader struct {
	in  midi.In
	key modular.V
}

func (r *midiKeyReader) Read(vs []modular.V) (n int, err error) {
	for i := range vs {
		vs[i] = r.key
	}
	return len(vs), nil
}

func (r *midiKeyReader) Close() error {
	if err := r.in.StopListening(); err != nil {
		return err
	}
	return r.in.Close()
}

type midiGateReader struct {
	in   midi.In
	gate modular.V
}

func (r *midiGateReader) Read(vs []modular.V) (n int, err error) {
	for i := range vs {
		vs[i] = r.gate
	}
	return len(vs), nil
}

func (r *midiGateReader) Close() error {
	if err := r.in.StopListening(); err != nil {
		return err
	}
	return r.in.Close()
}

func (i *Interface) Vel() modular.Reader {
	panic("midi.Interface.Vel: not implemented")
}

func (i *Interface) SparseKey() modular.SparseReader {
	panic("midi.Interface.SparseKey: not implemented")
}

func (i *Interface) SparseGate() modular.SparseReader {
	panic("midi.Interface.SparseGate: not implemented")
}

func (i *Interface) SparseVel() modular.Reader {
	panic("midi.Interface.SparseVel: not implemented")
}
