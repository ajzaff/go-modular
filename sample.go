package sample

// Sample represents a singular stereo audio sample.
//
// A complex number provides easy, native stereo channels.
type Sample complex128

// Store the sample value v into s.
func (s *Sample) Store(v complex128) {
	*s = Sample(v)
}

// StoreLeft stores the left channel into s.
func (s *Sample) StoreLeft(left float64) {
	s.Store(complex(left, s.Right()))
}

// StoreRight stores the right channel into s.
func (s *Sample) StoreRight(right float64) {
	s.Store(complex(s.Left(), right))
}

// Clear the sample.
func (s *Sample) Clear() {
	s.Store(0)
}

// Left returns the left channel of s.
func (s Sample) Left() float64 {
	return real(s)
}

// Right returns the right channel of s.
func (s Sample) Right() float64 {
	return imag(s)
}

// Reader is an interface for sample outputs.
type Reader interface {
	// Read accepts a sample buffer to read into.
	//
	// Implementations should return the number values read and whether
	// EOF has been reached (using standard io errors).
	Read(vs []Sample) (n int, err error)
}

// Writer is an interface for sample inputs.
type Writer interface {
	// Write accepts a sample buffer to write into.
	//
	// Implementations should return the number values written and whether
	// EOF has been reached (using standard io errors).
	Write(vs []Sample) (n int, err error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// Processor is an interface for signal processors.
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
