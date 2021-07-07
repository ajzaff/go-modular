package modio

import (
	"github.com/ajzaff/go-modular"
	"github.com/mjibson/go-dsp/fft"
)

type FFT struct {
	buf []complex128

	blockSize  int
	sampleRate int
}

func (x *FFT) Reset() {
	x.buf = nil
}

func (x *FFT) SetConfig(cfg *modular.Config) {
	x.blockSize = cfg.BufferSize
	x.sampleRate = cfg.SampleRate
}

// Compute FFT(b).
func (x *FFT) Compute(b []float32) []complex128 {
	p := make([]complex128, len(b))
	for i, v := range b {
		p[i] = complex(float64(v), 0)
	}
	return fft.FFT(p)
}

// Store FFT(b).
func (x *FFT) StoreFFT(b []float32) {
	if x.buf != nil {
		return
	}
	x.buf = x.Compute(b)
}

// Receive IFFT(x.buf) into b.
func (x *FFT) Process(b []float32) {
	for i, v := range fft.IFFT(x.buf) {
		b[i] = float32(real(v))
	}
}

// Returns FFT(b)_i.
func (x *FFT) Get(i int) complex128 {
	return x.buf[i]
}

// Update FFT(b)_i = v.
func (x *FFT) Update(i int, v complex128) {
	x.buf[i] = v
}

// Update all FFT(b)_i by applying f.
func (x *FFT) UpdateAll(f func(i int, v complex128) complex128) {
	for i, v := range x.buf {
		x.buf[i] = f(i, v)
	}
}

func CopyToReal32(dst []float32, src []complex128) {
	if len(dst) != len(src) {
		panic("modio.ToReal32: slices must have equal length")
	}
	for i := range dst {
		dst[i] = float32(real(src[i]))
	}
}

func CopyToComplex128(dst []complex128, src []float32) {
	if len(dst) != len(src) {
		panic("modio.ToComplex128: slices must have equal length")
	}
	for i := range dst {
		dst[i] = complex(float64(src[i]), 0)
	}
}
