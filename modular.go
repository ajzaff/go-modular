package modular

// V is a singular mono audio sample.
//
// Short for "voltage".
type V float64

// Store the sample value x into v.
func (v *V) Store(x float64) { *v = V(x) }

// Clear the sample.
func (v *V) Clear() { v.Store(0) }

// Context encapsulates modular synth context.
type Context struct {
	SampleRate int
	BufferSize int
	Driver
}

// Driver is an interface for synth drivers.
//
// The driver is usable after a call to Init.
// Send sends the entirety of the input voltage
// to output channel ch.
type Driver interface {
	// Init initializes the driver based on the context ctx.
	Init(ctx *Context)

	// Send sends audio data to the output.
	Send(ch int, in <-chan V) (n int64, err error)
}

// ContextOption allows configuration
type ContextOption func(*Context)

// WithSampleRate sets the context sample rate.
func WithSampleRate(sampleRate int) ContextOption {
	return func(ctx *Context) {
		ctx.SampleRate = sampleRate
	}
}

// WithBufferSize sets the context buffer size.
func WithBufferSize(bufferSize int) ContextOption {
	return func(ctx *Context) {
		ctx.BufferSize = bufferSize
	}
}

// NewContext returns a new modular context with default options.
//
// It calls the driver's Init method.
func NewContext(drv Driver, opts ...ContextOption) *Context {
	ctx := &Context{
		SampleRate: 44100,
		BufferSize: 44100,
		Driver:     drv,
	}
	for _, f := range opts {
		f(ctx)
	}
	drv.Init(ctx)
	return ctx
}
