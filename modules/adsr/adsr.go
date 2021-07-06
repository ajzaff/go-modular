package adsr

import (
	"math"
	"time"

	"github.com/ajzaff/go-modular"
)

// ADSR is a basic ADSR envelope generator.
type ADSR struct {
	a, d time.Duration
	s    float32
	r    time.Duration

	phase   phase
	begin   int
	p       int
	end     int
	sustain struct {
		set bool
		pos int
	}
	gate func() bool

	sampleRate int
}

func New(a time.Duration, d time.Duration, s float32, r time.Duration) *ADSR {
	adsr := &ADSR{
		a:          a,
		d:          d,
		s:          s,
		r:          r,
		sampleRate: 44100,
	}
	adsr.Reset()
	return adsr
}

func (a *ADSR) SetConfig(cfg *modular.Config) error {
	a.sampleRate = cfg.SampleRate
	a.Reset()
	return nil
}

type phase int

const (
	attack phase = iota
	decay
	sustain
	release
)

func (a *ADSR) samples(d time.Duration) int {
	return int(math.Round(float64(a.sampleRate) * float64(d.Seconds())))
}

// Reset manually resets the ADSR to the attack phase.
//
// Use SetGate to automate the reset.
func (a *ADSR) Reset() {
	a.phase = attack
	a.begin = 0
	a.p = 0
	a.end = a.samples(a.a)
}

// ResetSustain clears the fixed sustain duration.
func (a *ADSR) ResetSustain() {
	a.sustain = struct {
		set bool
		pos int
	}{}
}

// SetSustain optionally sets a fixed duration for the sustain phase.
//
// The duration d should not include attack and decay.
func (a *ADSR) SetSustain(d time.Duration) {
	a.sustain = struct {
		set bool
		pos int
	}{true, a.samples(a.a + a.d + d)}
}

// ResetGate unsets the automatic gate trigger.
func (a *ADSR) ResetGate() {
	a.gate = nil
}

// SetGate sets the automatic gate trigger.
//
// 	gate will be called once per sample.
// 	`gate() == true` resets the ADSR.
func (a *ADSR) SetGate(gate func() bool) {
	a.gate = gate
}

// Release releases the note now.
func (a *ADSR) Release() {
	if a.phase == sustain {
		a.releaseNow()
	}
}

func (a *ADSR) releaseNow() {
	a.phase = release
	a.begin = a.p
	a.end = a.p + a.samples(a.r)
}

// Envelope returns the next envelope amplitude.
//
// Envelope calls mutate the ADSR.
func (a *ADSR) Envelope() float32 {
	switch a.phase {
	case attack:
		if a.p >= a.end {
			a.phase = decay
			a.begin = a.p
			a.end = a.p + a.samples(a.d)
			return 1
		}
		defer func() { a.p++ }()
		// TODO: Configurable attackRamp func.
		return float32(a.p-a.begin) / float32(a.end-a.begin)
	case decay:
		if a.p >= a.end {
			a.phase = sustain
			a.begin = -1
			a.end = -1
			return a.s
		}
		defer func() { a.p++ }()
		return 1 - a.s*float32(a.p-a.begin)/float32(a.end-a.begin)
	case sustain:
		if a.sustain.set && a.p >= a.sustain.pos {
			a.releaseNow()
			return a.s
		}
		a.p++
		return a.s
	case release:
		if a.p >= a.end {
			return 0
		}
		defer func() { a.p++ }()
		// TODO: Configurable release ramp.
		return a.s - a.s*float32(a.p-a.begin)/float32(a.end-a.begin)
	default:
		panic("ADSR.Envelope: impossible state reached")
	}
}

// Process convolves the block with the ADSR envelope.
func (a *ADSR) Process(b []float32) {
	for i, v := range b {
		if a.gate != nil && a.gate() {
			a.Reset()
		}
		b[i] = v * a.Envelope()
	}
}
