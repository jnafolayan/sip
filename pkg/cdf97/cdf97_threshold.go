package cdf97

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

func (cdf *CDF97Wavelet) HardThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	return s.HardThreshold(0, 0, threshold)
}

func (cdf *CDF97Wavelet) SoftThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	return s.SoftThreshold(0, 0, threshold)
}
