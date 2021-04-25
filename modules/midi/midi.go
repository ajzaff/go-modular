package midi

import (
	"context"
	"fmt"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/portmididrv"
)

// Interface returns midi CVs for the stream of midi messages
// on the single midi input channel (i, ch).
//
// Interface is unbuffered to minimize trigger latency.
func Interface(ctx context.Context, i, ch uint8) (gate, key, vel control.CV, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("midi.Interface: %v", err)
		}
	}()

	drv, err := portmididrv.New()
	if err != nil {
		return nil, nil, nil, err
	}

	ins, err := drv.Ins()
	if len(ins) <= int(i) {
		return nil, nil, nil, err
	}

	in := ins[i]
	if err := in.Open(); err != nil {
		return nil, nil, nil, err
	}

	rd := reader.New(reader.NoLogger())

	gateCh := make(chan modular.V)
	keyCh := make(chan modular.V)
	velCh := make(chan modular.V)

	rd.Channel.NoteOn = func(p *reader.Position, channel, key, vel uint8) {
		if channel != ch {
			return
		}
		gateCh <- 1
		keyCh <- modular.V(key)
		velCh <- modular.V(vel) / 127
	}

	rd.Channel.NoteOff = func(p *reader.Position, channel, key, vel uint8) {
		if channel != ch {
			return
		}
		gateCh <- 0
		velCh <- 0
	}

	go func() {
		defer func() {
			close(gateCh)
			close(keyCh)
			close(velCh)
		}()
		if err := rd.ListenTo(in); err != nil {
			panic(fmt.Errorf("midi.Interface: %v", err))
		}
	}()

	return gateCh, keyCh, velCh, nil
}
