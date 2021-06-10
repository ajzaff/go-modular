package modio

// Repeat returns a new sample with b repeated count times.
//
// It panics if count is negative or if
// the result of (len(b) * count) overflows.
func Repeat(b []float32, count int) []float32 {
	if count == 0 {
		return []float32{}
	}
	if count < 0 {
		panic("modio.Repeat: negative Repeat count")
	} else if len(b)*count/count != len(b) {
		panic("modio.Repeat: count causes overflow")
	}
	nb := make([]float32, len(b)*count)
	bp := copy(nb, b)
	for bp < len(nb) {
		copy(nb[bp:], nb[:bp])
		bp *= 2
	}
	return nb
}
