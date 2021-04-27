package modular

import (
	"context"
	"errors"
)

// V is a mono audio sample.
//
// Short for voltage.
type V float64

// Store the sample value x into v.
func (v *V) Store(x float64) { *v = V(x) }

// Clear the sample.
func (v *V) Clear() { v.Store(0) }

var (
	sampleRateKey       int
	bufferSizeKey       int
	driverBufferSizeKey int
	driverKey           int
)

// Driver is an interface for raw audio output drivers.
//
// The driver is usable after a call to Init.
// Send sends the entirety of the input voltage
// to output channel ch.
type Driver interface {
	// Init initializes the driver based on the context ctx.
	Init(ctx context.Context)

	// InitContext initializes the driver based on the context ctx.
	InitContext(ctx *Context)

	// Send sends audio data to the output.
	Send(ch int, in <-chan V) (n int64, err error)

	// SendReader sends audio data to the output from r.
	SendReader(ch int, r Reader) (n int64, err error)
}

// New returns a new modular context from ctx with default options.
//
// It calls the driver's Init method.
func New(ctx context.Context, drv Driver) context.Context {
	ctx = context.WithValue(&modularContext{ctx}, &driverKey, drv)
	drv.Init(ctx)
	return ctx
}

type modularContext struct{ context.Context }

func (ctx *modularContext) Value(key interface{}) interface{} {
	switch key {
	case &sampleRateKey:
		return SampleRate(ctx.Context)
	case &bufferSizeKey:
		return BufferSize(ctx.Context)
	case &driverBufferSizeKey:
		return DriverBufferSize(ctx.Context)
	default:
		return ctx.Context.Value(key)
	}
}

// WithSampleRate sets the context sample rate.
func WithSampleRate(ctx context.Context, sampleRate int) context.Context {
	return context.WithValue(ctx, &sampleRateKey, sampleRate)
}

// WithBufferSize sets the context buffer size.
func WithBufferSize(ctx context.Context, bufferSize int) context.Context {
	return context.WithValue(ctx, &bufferSizeKey, bufferSize)
}

// WithDriverBufferSize sets the driver context buffer size.
func WithDriverBufferSize(ctx context.Context, bufferSize int) context.Context {
	return context.WithValue(ctx, &driverBufferSizeKey, bufferSize)
}

// SampleRate gets the sample rate for ctx.
func SampleRate(ctx context.Context) int {
	if v := ctx.Value(&sampleRateKey); v != nil {
		return v.(int)
	}
	return 44100
}

// BufferSize gets the buffer size for ctx.
func BufferSize(ctx context.Context) int {
	if v := ctx.Value(&bufferSizeKey); v != nil {
		return v.(int)
	}
	return 0
}

// DriverBufferSize gets the driver buffer size for ctx.
func DriverBufferSize(ctx context.Context) int {
	if v := ctx.Value(&driverBufferSizeKey); v != nil {
		return v.(int)
	}
	return 44100
}

// Send sends the input signal over the channel ch using the driver in context ctx.
func Send(ctx context.Context, ch int, in <-chan V) (n int64, err error) {
	if drv := ctx.Value(&driverKey); drv != nil {
		return drv.(Driver).Send(ch, in)
	}
	return 0, errors.New("bad driver")
}
