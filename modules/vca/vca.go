package vca

import (
	"context"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
)

// VCA is a simple voltage controlled amplifier.
//
// CV in is the audio signal.
// CV a is the amplitude voltage usually sourced from an ADSR envelope generator.
func VCA(ctx context.Context, a control.CV, in <-chan modular.V) <-chan modular.V {
	ch := make(chan modular.V)
	go func() {
		for v := range in {
			ch <- <-a * v
		}
		close(ch)
	}()
	return ch
}
