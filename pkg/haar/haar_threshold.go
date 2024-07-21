package haar

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

func (hw *HaarWavelet) HardThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	return s.HardThreshold(0, 0, threshold)
}

func (hw *HaarWavelet) SoftThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	return s.SoftThreshold(0, 0, threshold)
}
