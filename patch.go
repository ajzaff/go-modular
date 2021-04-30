package modular

import "io"

// Patch propogates the audio patch through the signal path to dst.
// Patch returns a cancel function which stops the propogation on demand.
//
// Write on the last path element is not used.
// UpdateConfig is automatically called on all elements which support it.
// Close is automatically called on all elements after terminating.
func (m *Modular) Patch(dst Writer, path ...Processor) (cancel func() error) {
	var err error
	done := make(chan struct{})

	if x, ok := dst.(Module); ok {
		if err1 := x.SetConfig(m.cfg); err == nil {
			err = err1
		}
	}
	for _, p := range path {
		if x, ok := p.(Module); ok {
			if err1 := x.SetConfig(m.cfg); err == nil {
				err = err1
			}
		}
	}

	go func() {
		defer close(done)
		if len(path) == 0 {
			return
		}
	loop:
		for {
			for i := len(path) - 1; i > 0; i-- {
				if _, err = copyBuffer(path[i-1], LimitReader(path[i], 512), nil, done); err != nil {
					break loop
				}
			}
			if _, err = copyBuffer(dst, LimitReader(path[0], 512), nil, done); err != nil {
				break
			}
		}
	}()
	return func() error {
		defer func() {
			if closer, ok := dst.(io.Closer); ok {
				if err1 := closer.Close(); err == nil {
					err = err1
				}
			}
			for _, p := range path {
				if closer, ok := p.(io.Closer); ok {
					if err1 := closer.Close(); err == nil {
						err = err1
					}
				}
			}
		}()
		defer func() { recover() }()
		done <- struct{}{} // done can be closed; hence recover()
		return err
	}
}

func (m *Modular) PatchSparse(dst SparseWriter, src ...SparseReader) (cancel func()) {
	panic("modular.PatchSparse: not implemented")
}
