package math

type Func func(float64) float64

func (fn Func) Process(buf []float64) {
	n := len(buf) / 2
	for i := range buf[:n] {
		buf[i] = fn(buf[i+n])
	}
}
