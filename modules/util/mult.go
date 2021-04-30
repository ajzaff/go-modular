package util

import (
	"sync"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/modio"
)

// Mult copies the input audio signal n times.
type Mult struct {
	outputs []*multOutput
	mu      sync.RWMutex
}

func (m *Mult) Write(vs []modular.V) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, out := range m.outputs {
		if n, err = out.buf.Write(vs); err != nil {
			return n, err
		}
	}
	return n, err
}

// New creates a new mult output.
func (m *Mult) New() modular.Reader {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := &multOutput{mu: &m.mu}
	m.outputs = append(m.outputs, out)
	return out
}

type multOutput struct {
	buf modio.Buffer
	mu  *sync.RWMutex
}

func (m *multOutput) Read(vs []modular.V) (n int, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.buf.Read(vs)
}
