package midi

import (
	"fmt"

	"github.com/ajzaff/go-modular"
	"gitlab.com/gomidi/rtmididrv"
)

// Interface returns midi CVs for the stream of midi messages
// on the single midi input channel (i, ch).
//
// Interface is unbuffered to minimize trigger latency.
type Interface struct {
	drv *rtmididrv.Driver
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

	iface = &Interface{}
	return iface, nil
}

func (i *Interface) UpdateConfig(cfg *modular.Config) error {
	return nil
}

func SparseGate() modular.SparseReader {
	panic("midi.Interface.SparseGate: not implemented")
}

func Key() modular.Reader {
	panic("midi.Interface.Key: not implemented")
}

func SparseKey() modular.SparseReader {
	panic("midi.Interface.SparseKey: not implemented")
}

func Vel() modular.Reader {
	panic("midi.Interface.Vel: not implemented")
}

func SparseVel() modular.Reader {
	panic("midi.Interface.SparseVel: not implemented")
}
