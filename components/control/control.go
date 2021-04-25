package control

import (
	"context"
	"math"
	"sync"
	"sync/atomic"

	"github.com/ajzaff/go-modular"
)

// CV is a control voltage.
//
// The alias can be used to differentiate between audio and control voltages.
type CV <-chan modular.V

// Voltage returns a CV from a singluar value v.
//
// Equivalent to calling Func with a constant-yielding func.
func Voltage(ctx context.Context, v float64) CV {
	return Func(ctx, func() modular.V { return modular.V(v) })
}

// Func returns a variable voltage source from evaluating fn.
func Func(ctx context.Context, fn func() modular.V) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	go func() {
		for {
			ch <- fn()
		}
	}()
	return ch
}

// Latch takes a trigger CV and binds it to an continuous output.
//
// Latch can be used when the input CV has a trigger but a
// continuous output is desired.
//
//	cv1 := make(chan V)
//	go func() {
//		time.Sleep(Second)
//		cv1 <- 1 // triger start
//		time.Sleep(Second)
//		cv1 <- 2 // change value
//	}()
//	// <-cv1 // (after 1 second) 1
//	// <-cv1 // (after 2 seconds) 2
//	cv2 := Latch(cv1)
//	// <-cv2 // (after t<2 seconds) 1
//	// <-cv2 // (after t>=2 seconds) 2...
//
// Latch input should be unbuffered to minimize trigger latency.
func Latch(ctx context.Context, in CV) CV {
	ch := make(chan modular.V, modular.BufferSize(ctx))
	done := int32(0)

	var first sync.WaitGroup
	first.Add(1)

	var a uint64
	go func() {
		defer func() { atomic.StoreInt32(&done, 1) }()
		v, ok := <-in
		if !ok {
			return
		}
		first.Done()
		for {
			atomic.StoreUint64(&a, math.Float64bits(float64(v)))
			v, ok = <-in
			if !ok {
				break
			}
		}
	}()
	go func() {
		first.Wait()
		for atomic.LoadInt32(&done) == 0 {
			v := atomic.LoadUint64(&a)
			ch <- modular.V(math.Float64frombits(v))
		}
	}()
	return ch
}
