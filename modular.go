package modular

import (
	"context"
	"errors"
	"time"
)

// V is a singular mono audio sample.
//
// Short for voltage.
type V float64

// Store the sample value x into v.
func (v *V) Store(x float64) { *v = V(x) }

// Clear the sample.
func (v *V) Clear() { v.Store(0) }

var (
	sampleRateKey int
	bufferSizeKey int
	driverKey     int
)

// Driver is an interface for synth drivers.
//
// The driver is usable after a call to Init.
// Send sends the entirety of the input voltage
// to output channel ch.
type Driver interface {
	// Init initializes the driver based on the context ctx.
	Init(ctx context.Context)

	// Send sends audio data to the output.
	Send(ch int, in <-chan V) (n int64, err error)
}

// NewContext returns a new modular context with default options.
//
// It calls the driver's Init method.
func NewContext(drv Driver) context.Context {
	ctx := &modularContext{
		values: map[interface{}]interface{}{
			&sampleRateKey: 44100,
			&bufferSizeKey: 44100,
			&driverKey:     drv,
		},
	}
	drv.Init(ctx)
	return ctx
}

type modularContext struct {
	values map[interface{}]interface{}
}

func (ctx *modularContext) Deadline() (deadline time.Time, ok bool) { return }
func (ctx *modularContext) Done() <-chan struct{}                   { return nil }
func (ctx *modularContext) Err() error                              { return nil }
func (ctx *modularContext) Value(key interface{}) interface{}       { return ctx.values[key] }

// WithSampleRate sets the context sample rate.
func WithSampleRate(ctx context.Context, sampleRate int) context.Context {
	return context.WithValue(ctx, &sampleRateKey, sampleRate)
}

// WithBufferSize sets the context buffer size.
func WithBufferSize(ctx context.Context, bufferSize int) context.Context {
	return context.WithValue(ctx, &bufferSizeKey, bufferSize)
}

// SampleRate gets the sample rate for ctx.
func SampleRate(ctx context.Context) int {
	return ctx.Value(&sampleRateKey).(int)
}

// BufferSize gets the buffer size for ctx.
func BufferSize(ctx context.Context) int {
	return ctx.Value(&bufferSizeKey).(int)
}

// Send sends the input signal over the channel ch using the driver in context ctx.
func Send(ctx context.Context, ch int, in <-chan V) (n int64, err error) {
	if drv, ok := ctx.Value(&driverKey).(Driver); ok {
		return drv.Send(ch, in)
	}
	return 0, errors.New("bad driver")
}
