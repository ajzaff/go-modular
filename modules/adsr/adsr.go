package adsr

import (
	"time"

	"github.com/ajzaff/go-modular"
)

type triggerState int

const (
	stateAwait triggerState = iota
	stateAttack
	stateDecay
	stateSustain
	stateRelease
)

// Envelope is a basic ADSR envelope generator.
//
// Gate transitions shape the generated envelope:
//	0->1: Attack trigger.
//	1->1: Decay/Sustain.
//	1->0: Release trigger.
//	0->0: Off.
type Envelope struct {
	a, d       time.Duration
	s          float64
	r          time.Duration
	buf        []modular.V
	t          float64
	len        float64
	state      triggerState
	sampleRate int
}

func New(a time.Duration, d time.Duration, s float64, r time.Duration) *Envelope {
	return &Envelope{
		a,
		d,
		s,
		r,
		nil,
		0,
		0,
		stateAwait,
		44100,
	}
}

func (e *Envelope) UpdateConfig(cfg *modular.Config) error {
	e.sampleRate = cfg.SampleRate
	return nil
}

func (e *Envelope) Gate() modular.SparseWriter {
	panic("adsr.Envelope.Gate: not implemented")
}

func length(sampleRate int, d time.Duration) float64 {
	if d == 0 {
		return 0
	}
	return float64(sampleRate) * d.Seconds()
}

func (e *Envelope) Read(vs []modular.V) (n int, err error) {
	n = copy(vs, e.buf)
	for i, v := range vs[:n] {
		switch e.state {
		case stateAwait:
			if v > 0 {
				e.state = stateAttack
				e.t = 0
				e.len = length(e.sampleRate, e.a)
			}
			vs[i] = 0
		case stateAttack:
			if v <= 0 {
				e.state = stateRelease
				e.t = 0
				e.len = length(e.sampleRate, e.r)
				continue
			}
		case stateDecay:
			if v <= 0 {
				e.state = stateRelease
				e.t = 0
				e.len = length(e.sampleRate, e.r)
				continue
			}
		case stateSustain:
			if v <= 0 {
				e.state = stateRelease
				e.t = 0
				e.len = length(e.sampleRate, e.r)
				continue
			}
			vs[i] = modular.V(e.s)
		case stateRelease:
			if v > 0 {
				e.state = stateAttack
				e.t = 0
				e.len = length(e.sampleRate, e.a)
				continue
			}
			vs[i] = modular.V(e.s - e.s*e.t/e.len)
			e.t++
			if e.len < e.t+1 {
				e.state = stateAwait
				e.t = 0
				e.len = 0
			}
		}
	}
	return n, err
}
