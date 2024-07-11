package cdf97

type CDF97Wavelet struct {
	Level int
}

func (cdf *CDF97Wavelet) GetDecompositionLevel() int {
	return cdf.Level
}
