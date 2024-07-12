package cdf97

import "github.com/jnafolayan/sip/pkg/signal"

// Inspired by https://gist.github.com/CoderSherlock/834e9eb918eeb0dfee5f4550077f57f8
type CDF97Wavelet struct {
	Level int
}

func (cdf *CDF97Wavelet) GetDecompositionLevel() int {
	return cdf.Level
}

func transpose(s signal.Signal2D) signal.Signal2D {
	w, h := s.Size()
	t := signal.New(h, w)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			t[j][i] = s[i][j]
		}
	}
	return t
}
