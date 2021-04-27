package control

import (
	"math"
	"sync/atomic"

	"github.com/ajzaff/go-modular"
)

type Atomic uint64

func (a *Atomic) Store(v modular.V) {
	atomic.StoreUint64((*uint64)(a), math.Float64bits(float64(v)))
}

func (a *Atomic) Load() modular.V {
	return modular.V(math.Float64frombits(atomic.LoadUint64((*uint64)(a))))
}
