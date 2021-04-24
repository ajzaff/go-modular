package midi

import "math"

// Tuning is an interface for a chromatic scale.
type Tuning interface {
	// A4Hz returns the frequency of the note A4.
	A4Hz() float64
}

// StdTuning is an A440 tuning.
var StdTuning = stdTuning{}

type stdTuning struct{}

func (stdTuning) A4Hz() float64 {
	return 440
}

func clampKey(k int) uint8 {
	if k < 0 {
		k = 0
	}
	if k > 127 {
		k = 127
	}
	return uint8(k)
}

// Key returns the MIDI key value from tuning t and frequency f.
// A4 is MIDI key 69 for instance.
func Key(t Tuning, f float64) int {
	return int(69 + 12*math.Log2(f/t.A4Hz()))
}

// Pitch returns the pitch of the midi key in tuning t.
func Pitch(t Tuning, key int) float64 {
	return Tone(t, float64(key))
}

// Tone returns the tone of the fractional midi key in tuning t.
func Tone(t Tuning, key float64) float64 {
	return t.A4Hz() * math.Pow(2, (key-69)/12)
}

// Note constants starting at octave 0.
const (
	C = 12 + iota // C0
	Db
	D
	Eb
	E
	F
	Gb
	G
	Ab
	A
	Bb
	B
)

// Note returns the midi note in octave oct.
//
// Example:
//	A4 = Note(A, 4)
//
// Piano notes range from A0 (21) to C8 (108).
// Midi notes range from C-1 (0) to G9 (127).
func Note(note, oct int) int {
	return note + 12*oct
}
