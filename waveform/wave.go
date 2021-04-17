package waveform

import (
	"io"
	"math"

	"github.com/ajzaff/go-sample"
)

type waveform struct {
	fn  func(float64) float64
	buf *sample.Buffer
}

func (w *waveform) BlockSize() int { return 0 }

func (w *waveform) Write(vs []sample.Sample) (n int, err error) {
	return w.buf.Write(vs)
}

func (w *waveform) Read(vs []sample.Sample) (n int, err error) {
	n, err = w.buf.Read(vs)
	if err == io.EOF {
		err = nil
	}
	for i, v := range vs[:n] {
		vs[i].StoreLeft(w.fn(v.Left()))
	}
	return n, err
}

func Sine() sample.Processor {
	return &waveform{fn: func(v float64) float64 { return math.Sin(2 * math.Pi * v) }}
}

func Sinc() sample.Processor {
	return &waveform{fn: func(v float64) float64 { return math.Sin(math.Pi*v) / (math.Pi * v) }}
}

func Sawtooth() sample.Processor {
	return &waveform{fn: func(v float64) float64 { return 2 / math.Pi * math.Atan(math.Tan(math.Pi*v)) }}
}

func Triangle() sample.Processor {
	return &waveform{fn: func(v float64) float64 { return 2 / math.Pi * math.Asin(math.Sin(2*math.Pi*v)) }}
}
