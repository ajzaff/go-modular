package util

import (
	"io"

	"github.com/ajzaff/go-modular"
	"github.com/ajzaff/go-modular/sampleio"
)

func Reader(in <-chan modular.V) sampleio.Reader {
	return &reader{in}
}

type reader struct {
	in <-chan modular.V
}

func (r *reader) Read(vs []modular.V) (n int, err error) {
	i := 0
	for v := range r.in {
		vs[i] = v
		if i++; i == len(vs) {
			return i, nil
		}
	}
	return i, io.EOF
}
