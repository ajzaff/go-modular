package main

import (
	"math"

	"github.com/ajzaff/go-sample"
	"github.com/ajzaff/go-sample/outputs"
)

func main() {
	const sampleRate = 44100
	const freq = 440

	buf := make([]sample.Sample, 1<<11)
	oto := outputs.Oto(sampleRate)
	count := uint64(0)

	for {
		ct := count
		for i := range buf {
			v := math.Sin(2 * math.Pi * float64(ct) * freq / sampleRate)
			buf[i].Store(complex(v, v))
			ct++
		}

		n, _ := oto.Write(buf)
		count += uint64(n)
	}
}
