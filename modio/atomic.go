package modio

import (
	"math"
	"sync/atomic"

	"github.com/ajzaff/go-modular"
)

// Atomic uses atomic Store to package a sample into a uint64.
//
// Deprecated: Prefer using Sparse IO.
type Atomic uint64

func (a *Atomic) Store(v modular.V) {
	atomic.StoreUint64((*uint64)(a), math.Float64bits(float64(v)))
}

func (a *Atomic) Load() modular.V {
	return modular.V(math.Float64frombits(atomic.LoadUint64((*uint64)(a))))
}

func (a Atomic) LoadUnsafe() modular.V {
	return modular.V(math.Float64frombits(uint64(a)))
}
