package adsr

import (
	"context"
	"math"
	"time"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/components/control"
)

type triggerState int

const (
	stateAwait triggerState = iota
	stateAttack
	stateDecay
	stateSustain
	stateRelease
)

func length(sampleRate int, d time.Duration) float64 {
	if d == 0 {
		return 0
	}
	return float64(sampleRate) * d.Seconds()
}

// Envelope is a basic ADSR envelope generator.
//
// Gate transitions shape the generated envelope:
//	0->1: Attack trigger.
//	1->1: Decay/Sustain.
//	1->0: Release trigger.
//	0->0: Off.
func Envelope(ctx context.Context, a, d time.Duration, s float64, r time.Duration, gate control.CV) control.CV {
	ch := make(chan modular.V)
	go func() {
		var (
			t          float64
			len        float64
			sampleRate = modular.SampleRate(ctx)
			state      triggerState
		)
		for v := range gate {
			switch state {
			case stateAwait:
				if v > 0 {
					state = stateAttack
					t = 0
					len = length(sampleRate, a)
				}
				ch <- 0
			case stateAttack:
				if v <= 0 {
					state = stateRelease
					t = 0
					len = length(sampleRate, r)
					continue
				}
				t++
				ch <- modular.V(t / len)
				if len < t+1 {
					state = stateDecay
					t = 0
					len = length(sampleRate, d)
				}
			case stateDecay:
				if v <= 0 {
					state = stateRelease
					t = 0
					len = length(sampleRate, r)
					continue
				}
				ch <- modular.V(1 - (1-s)*t/len)
				t++
				if len < t+1 {
					state = stateSustain
					t = 0
					len = math.Inf(+1)
				}
			case stateSustain:
				if v <= 0 {
					state = stateRelease
					t = 0
					len = length(sampleRate, r)
					continue
				}
				ch <- modular.V(s)
			case stateRelease:
				if v > 0 {
					state = stateAttack
					t = 0
					len = length(sampleRate, a)
					continue
				}
				ch <- modular.V(s - s*t/len)
				t++
				if len < t+1 {
					state = stateAwait
					t = 0
					len = 0
				}
			}
		}
	}()
	return ch
}
