package midi

import (
	"fmt"
	"math"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/sampleio"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midiwriter"
	"gitlab.com/gomidi/portmididrv"
)

func MidiWriter(driver midi.Writer) sampleio.Writer {
	return &midiWriter{driver}
}

func PortMidiDriver(n int) midi.Writer {
	drv, err := portmididrv.New()
	if err != nil {
		panic(err)
	}
	defer drv.Close()
	outs, err := drv.Outs()
	if err != nil {
		panic(err)
	}
	if len(outs) <= n {
		panic("midi.PortMidiWriter: missing output")
	}
	out := outs[n]
	if err := out.Open(); err != nil {
		panic(fmt.Errorf("midi.PortMidiWriter: %w", err))
	}
	return midiwriter.New(out, midiwriter.NoRunningStatus())
}

type midiWriter struct {
	driver midi.Writer
}

func (w *midiWriter) Write(vs []modular.V) (n int, err error) {
	for _, v := range vs {
		writeOne(w.driver, float64(v))
	}
	return n, nil
}

func writeOne(w midi.Writer, x float64) {
	switch x := math.Float64bits(x); MessageType(x) {
	case MessageNoteOff:
		w.Write(channel.Channel0.NoteOff(uint8((x & 0xff000000000000) >> 48)))
	case MessageNoteOn:
		w.Write(channel.Channel0.NoteOn(uint8((x&0xff000000000000)>>48), uint8((x&0xff0000000000)>>40)))
	case MessageRealtimeReset:
		w.Write(realtime.Reset)
	}
}
