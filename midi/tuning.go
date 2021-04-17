package midi

import "math"

type Tuning interface {
	A4Hz() float64
}

var StdTuning = stdTuning{}

type stdTuning struct{}

func (stdTuning) A4Hz() float64 {
	return 440
}

// Key returns the clamped MIDI key value.
// A4 is MIDI key 69 for instance.
func Key(v int) uint8 {
	if v < 0 {
		v = 0
	}
	if v > 127 {
		v = 127
	}
	return uint8(v)
}

// Pitch uses tuning t to get the frequency of the note hs half steps away from A4.
func Pitch(t Tuning, hs int) float64 {
	return Tone(t, float64(hs))
}

const twelfthRoot2 = 1.0594630943592953

// Tone uses tuning t to get the frequency of the tone hs fractional half steps away from A4.
func Tone(t Tuning, hs float64) float64 {
	return t.A4Hz() * math.Pow(twelfthRoot2, hs)
}

const (
	A4b = 69 + iota - 1
	A4
	B4b
	B4
	C4b
	C4
	D4b
	D4
	E4b
	E4
	F4b
	F4
	G4b
	G4
)
