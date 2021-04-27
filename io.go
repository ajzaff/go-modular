package modular

// Reader is an interface for sample outputs.
type Reader interface {
	// Read accepts a sample buffer to read into.
	//
	// Implementations should return the number values read and whether
	// EOF has been reached (using standard io errors).
	Read(vs []V) (n int, err error)
}

// Writer is an interface for sample inputs.
type Writer interface {
	// Write accepts a sample buffer to write into.
	//
	// Implementations should return the number values written and whether
	// an error has been reached (using standard io errors).
	Write(vs []V) (n int, err error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// Processor is an interface for block sample processors.
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
