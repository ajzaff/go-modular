package modular

// V is a mono PCM audio sample.
//
// Short for voltage.
type V float32

// Config provides synth configuration and options.
type Config struct {
	// SampleRate configures how many samples are played per second.
	//
	// Defaults to 44.1k.
	SampleRate int
	// BufferSize configures the size of module buffers.
	//
	// Defaults to 44.1k.
	BufferSize int
	// DriverBufferSize configures the size of driver buffers.
	//
	// Defaults to 44.1k.
	//
	// Note that drivers require a larger buffer to be performant,
	// but sound might not begin playing until the buffer is full.
	DriverBufferSize int
	// SampleSize configures the sample size in number of samples.
	//ProductLen
	// The default is 512 (about 12ms at the standard sample rate).
	SampleSize int
}

// NewConfig returns a new Config with the default options.
func NewConfig() *Config {
	return &Config{
		SampleRate:       44100,
		BufferSize:       44100,
		DriverBufferSize: 44100,
		SampleSize:       512, // ~12ms
	}
}

// BlockSize is an interface for modules requiring a minimum number of samples.
//
// For instance, block filters require a ceterain amount of samples to work.
//
// By default the blocksize is determined by the synth engine and assumed to be
// 1:1 for input:output.
type BlockSize interface {
	// BlockSize returns the minimum input and output block sizes.
	//
	// Module.Process must be called in multiples of BlockSize.
	BlockSize() (in, out int)
}

// Shape is an interface for modules with custom inputs and outputs.
//
// By default there is assumed to be a single input and output.
type Shape interface {
	// Shape returns the number of inputs and outputs for this Module.
	//
	// Special cases are in=0: a source Module, out=0: a sink Module.
	// in=out=0: Nil Module.
	Shape() (in, out int)
}

// Module is the interface for types that expect modular configuration.
//
// This is also expected on drivers to allow for updating configuration on demand.
type Module interface {
	// SetConfig updates the module configuration to cfg.
	//
	// Modules should support calling UpdateConfig arbitrarily many times.
	// An error may be returned either from Close on the old driver or using
	// unsupported config.
	//
	// If no Config is provided when creating a module, the driver may not be
	// initialized until after the first call to UpdateConfig which happens
	// within the New constructor.
	SetConfig(cfg *Config) error

	// Processes the audio buffer and returns the resulting data slice.
	//
	// The buf is structured particularly with space preallocated for the
	// inputs and outputs; having inputs concatenated after outputs:
	//
	//	[ out1..out2..out3..out_n | in1..in2..in3..in_n ]
	//
	// The optional Shape interface method controls the number of
	// inputs and outputs which must not change during the lifetime of
	// a Module.
	//
	// The optional Blocksize interface method controls the minimum
	// number of samples for each in_i and out_i. Otherwise, this value
	// is determined by the synth engine.
	//
	// Module must not retain buf.
	Process(buf []V)
}

type Modular struct {
	cfg *Config
}

// New creates a modular with default config.
func New() (*Modular, error) {
	return NewWithConfig(nil)
}

// New creates a modular with the given config.
//
// If cfg is nil the defaults will be used.
func NewWithConfig(cfg *Config) (*Modular, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	return &Modular{cfg: cfg}, nil
}
