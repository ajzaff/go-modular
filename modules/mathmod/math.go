package mathmod

type Func func(float32) float32

func (fn Func) Process(b []float32) {
	for i, v := range b {
		b[i] = fn(v)
	}
}

type Func2 func(i int, v float32) float32

func (fn Func2) Process(b []float32) {
	for i, v := range b {
		b[i] = fn(i, v)
	}
}
