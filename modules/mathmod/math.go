package mathmod

import "github.com/ajzaff/go-modular"

type Func func(float32) float32

func (fn Func) Process(b modular.Block) {
	for i, v := range b.Buf {
		b.Buf[i] = fn(v)
	}
}

type Func2 func(int, float32) float32

func (fn2 Func2) Process(b modular.Block) {
	for i, v := range b.Buf {
		b.Buf[i] = fn2(b.Pos+i, v)
	}
}
