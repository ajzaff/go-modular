package outputs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ajzaff/go-sample"
	"github.com/hajimehoshi/oto"
)

// Oto outputs to the speaker using the Oto mechanism.
// Each call to Oto starts a new stateful speaker context.
// This means you should not continously invoke Oto but rather
// Wrap the call in a Compose, Forever, or Repeat statement.
func Oto(sampleRate int) sample.Writer {
	return OtoContext(sampleRate, 2, 2, 2*sampleRate)
}

func OtoContext(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes int) sample.Writer {
	ctx, err := oto.NewContext(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(fmt.Errorf("outputs.OtoContext: %w", err))
	}
	var oto otoContext
	oto.player = ctx.NewPlayer()
	return &oto
}

type otoContext struct {
	player *oto.Player
}

func (o *otoContext) Write(vs []sample.Sample) (n int, err error) {
	var buf bytes.Buffer
	buf.Grow(4 * len(vs))
	for _, s := range vs {
		binary.Write(&buf, binary.LittleEndian, convert(s.Left()))
		binary.Write(&buf, binary.LittleEndian, convert(s.Right()))
	}
	written, err := io.Copy(o.player, &buf)
	return int(written), err
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
