package midi

import "math"

// Tuning is an interface for a chromatic scale.
type Tuning interface {
	// A4Hz returns the frequency of the note A4.
	A4Hz() float32
}

// StdTuning is an A440 tuning.
var StdTuning = stdTuning{}

type stdTuning struct{}

func (stdTuning) A4Hz() float32 {
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
func Key(t Tuning, f float32) int {
	return int(69 + 12*math.Log2(float64(f)/float64(t.A4Hz())))
}

// Pitch returns the pitch of the midi key in tuning t.
func Pitch(t Tuning, key int) float32 {
	return Tone(t, float32(key))
}

// Tone returns the tone of the fractional midi key in tuning t.
func Tone(t Tuning, key float32) float32 {
	return t.A4Hz() * float32(math.Pow(2, (float64(key)-69)/12))
}

// Note constants starting at octave 0.
const (
	C  = 12 + iota // C0
	Db             // Db0
	D              // D0
	Eb             // Eb0
	E              // E0
	F              // F0
	Gb             // Gb0
	G              // G0
	Ab             // Ab0
	A              // A0
	Bb             // Bb0
	B              // B0
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
