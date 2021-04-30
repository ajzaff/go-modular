package modular

import (
	"fmt"
	"io"
	"sync"
)

// V is a mono audio sample.
//
// Short for voltage.
type V float64

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
	//
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
}

// Driver is an interface for raw audio output drivers.
//
// Send sends the entirety of the input voltage
// to output channel ch.
type Driver interface {
	Module
	io.Closer

	// Send returns a new writer for audio channel ch.
	Send(ch int) WriteCloser
}

type Modular struct {
	Driver
	cfg *Config
}

// New creates a modular with default config.
func New(drv Driver) (*Modular, error) {
	return NewWithConfig(drv, nil)
}

// New creates a modular with the given config.
//
// If cfg is nil the defaults will be used.
func NewWithConfig(drv Driver, cfg *Config) (*Modular, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	if err := drv.SetConfig(cfg); err != nil {
		return nil, fmt.Errorf("modular.New: %w", err)
	}
	return &Modular{Driver: drv, cfg: cfg}, nil
}

// Reader is an interface for sample source modules.
type Reader interface {
	// Read accepts a sample buffer to read into.
	//
	// Implementations should return the number values read and whether
	// EOF has been reached (using standard io errors).
	Read(vs []V) (n int, err error)
}

// LimitReader reads up to n samples from r.
func LimitReader(r Reader, n int64) Reader {
	return &limitReader{r, n}
}

type limitReader struct {
	r Reader
	n int64
}

func (r *limitReader) Read(vs []V) (n int, err error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(vs)) > r.n {
		vs = vs[:r.n]
	}
	n, err = r.r.Read(vs)
	r.n -= int64(n)
	return
}

// Writer is an interface for sample inputs.
type Writer interface {
	// Write accepts a sample buffer to write into.
	//
	// Implementations should return the number values written and whether
	// an error has been reached (using standard io errors).
	Write(vs []V) (n int, err error)
}

// Processor is an interface for sample processing modules.
type Processor interface {
	Reader
	Writer

	// BlockSize returns the suggested block size for the processor.
	//
	// Read and Write can be called in smaller chunks, but work is
	// optimized on chunks of BlockSize samples.
	// BlockSize should return 0 when there is no preferred block size.
	BlockSize() int
}

// NopProcessor wraps r with a nop Write and zero BlockSize.
//
// Useful for passing a Reader as the last element of Patch.
func NopProcessor(r Reader) Processor { return &nopProcessor{r} }

type nopProcessor struct{ Reader }

func (p *nopProcessor) Write(vs []V) (n int, err error) { return len(vs), nil }
func (p *nopProcessor) BlockSize() int                  { return 0 }

// WriteCloser is an interface for writers with a Close method.
//
// Usually reseved for driver channels.
type WriteCloser interface {
	Writer
	io.Closer
}

// WriteNotifier returns a new writer and recv-only channel notifying
// upon the first write to w.
//
// This is useful for blocking operations which expect w to have been
// written.
func WriteNotifier(w Writer) (out Writer, firstWrite <-chan struct{}) {
	firstWrite = make(chan struct{})
	return &writeNotifier{w: w}, firstWrite
}

type writeNotifier struct {
	w      Writer
	notify chan struct{}
	once   sync.Once
}

func (w *writeNotifier) Write(vs []V) (n int, err error) {
	w.once.Do(func() {
		w.notify <- struct{}{}
		close(w.notify)
	})
	return w.w.Write(vs)
}

// SparseV is used with the Sparse Reader and Writer interfaces.
// It provides an efficient solution for implementing sparse CV
// signals such as gates.
type SparseV struct {
	// Idx is the relative index of the sample.
	// Negative values for Idx are never valid.
	Idx int
	// V is the control voltage sample value.
	V
}

// SparseReader is an interface for reading sparse CV samples.
//
// It saves bandwidth when encoding values that change infrequently.
type SparseReader interface {
	// SparseRead populates vs with control voltages
	// and returns the number of CVs read and error if any.
	//
	// Implementations can assume vs is ordered by ascending Idx.
	// See SparseV for explaination of sparse values.
	ReadSparse(vs []SparseV) (n int, err error)
}

// SparseWriter is an interface for writing sparse CV samples.
//
// It saves bandwidth when encoding values that change infrequently.
type SparseWriter interface {
	// SparseWrite writes vs to the writer and returns
	// the number of CVs written and error if any.
	//
	// Implementations can assume vs is ordered by ascending Idx.
	// See SparseV for explaination of sparse values.
	WriteSparse(vs []SparseV) (n int, err error)
}

// WrapReader wraps r as a SparseReader which reads a V on each value change.
//
// Depending on the density of the reader r there may be no benefits to
// converting to a SparseReader.
func WrapReader(r Reader) SparseReader {
	return &wrappedReader{r: r}
}

type wrappedReader struct {
	r   Reader
	p   int
	v   V
	buf []V
}

func (r *wrappedReader) ReadSparse(vs []SparseV) (n int, err error) {
	if r.buf == nil {
		r.buf = make([]V, 512)
	}
	n1, err := r.r.Read(r.buf)
	for _, v := range r.buf[:n1] {
		if v != r.v {
			r.v = v
			vs[n].Idx = r.p
			n++
		}
		r.p++
		if len(vs) <= n {
			break
		}
	}
	return n, err
}

// WrapWriter wraps w as a SparseWriter which accepts sparse values.
//
// Depending on the density of the Writer w there may be no benefits to
// converting to a SparseReader.
func WrapWriter(w Writer) SparseWriter {
	return &wrappedWriter{w: w}
}

type wrappedWriter struct {
	w   Writer
	p   int
	v   V
	buf []V
}

func (w *wrappedWriter) WriteSparse(vs []SparseV) (n int, err error) {
	if w.buf == nil {
		w.buf = make([]V, 512)
	}
	for i := range w.buf {
		if w.p == vs[0].Idx {
			w.v = vs[0].V
		}
		w.buf[i] = w.v
		w.p++
	}
	n, err = w.w.Write(w.buf)
	return n, err
}
