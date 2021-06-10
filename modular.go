package modular

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
	// The default is 512 (about 12ms at the standard sample rate).
	SampleSize int

	// Phase configures the phase shift value for supported modules.
	Phase float32
}

// New returns a new modular config with default values.
func New() *Config {
	return &Config{
		SampleRate:       44100,
		BufferSize:       44100,
		DriverBufferSize: 44100,
		SampleSize:       512, // ~12ms
	}
}

// OutputScaled is an optional interface for modules with custom output scale.
type OutputScaled interface {
	// OutputScaled returns the output scale for this Module.
	//
	// The number of times larger to scale the output.
	OutputScaled() float32
}

// Block is a mono PCM audio sample.
//
// Pos is the absolute position of the sample.
type Block struct {
	Buf []float32
	Pos int
}

// Processor is an interface for block processors.
type Processor interface {
	// Process processes the audio block in place.
	//
	// The optional OutputScaled interface method controls the number of
	// extra samples in the buffer. The number of extra samples will be
	// the larger of OutputScaled-1 and 0.
	//
	// The processor must not retain Buf.
	Process(b Block)
}

// Module is an interface for configurable block processors.
type Module interface {
	Processor

	// SetConfig updates the module configuration to cfg.
	//
	// Modules should expect UpdateConfig once for initialization and allow
	// arbitrarily many times after.
	SetConfig(cfg *Config) error
}
