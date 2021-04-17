package midi

type Time struct {
	// Number of samples per second.
	SampleRate int

	// BPM in number of D1_4 beats per minute.
	BPM int
}

// Duration is a value equal to number of beats per second at 60 BPM.
type Duration float64

// Common duration constants.
const (
	D1_4  Duration = 1
	D1_8  Duration = D1_4 / 2
	D1_16 Duration = D1_8 / 2
	D1_32 Duration = D1_16 / 2

	D1_4d  Duration = 1.5
	D1_8d  Duration = D1_4d / 2
	D1_16d Duration = D1_8d / 2
	D1_32d Duration = D1_16d / 2
)

// N converts a duration d to a number in raw samples.
func (t Time) N(d Duration) int {
	return int(60 / Duration(t.BPM) * Duration(t.SampleRate) * d)
}
