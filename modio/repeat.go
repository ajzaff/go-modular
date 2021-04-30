package modio

import "github.com/ajzaff/go-modular"

// Repeat returns a new sample with b repeated count times.
//
// It panics if count is negative or if
// the result of (len(b) * count) overflows.
func Repeat(b []modular.V, count int) []modular.V {
	if count == 0 {
		return []modular.V{}
	}
	if count < 0 {
		panic("modio.Repeat: negative Repeat count")
	} else if len(b)*count/count != len(b) {
		panic("modio.Repeat: count causes overflow")
	}
	nb := make([]modular.V, len(b)*count)
	bp := copy(nb, b)
	for bp < len(nb) {
		copy(nb[bp:], nb[:bp])
		bp *= 2
	}
	return nb
}
