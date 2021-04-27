package modular

import "sync"

// Sample is a short mono audio sample.
type Sample []float64

var samplePool sync.Pool

func init() {
	samplePool.New = func() interface{} {
		s := make(Sample, 4096)
		return &s
	}
}

func GetSample() Sample {
	return *(samplePool.Get().(*Sample))
}

func FreeSample(s *Sample) {
	samplePool.Put(s)
}
