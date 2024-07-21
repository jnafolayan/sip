package haar

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

func (hw *HaarWavelet) HardThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	width, height := s.Size()
	offsetX := width / (1 << hw.Level)
	offsetY := height / (1 << hw.Level)

	return s.HardThreshold(offsetX, offsetY, threshold)
}

func (hw *HaarWavelet) SoftThreshold(s signal.Signal2D, threshold int) signal.Signal2D {
	width, height := s.Size()
	offsetX := width / (1 << hw.Level)
	offsetY := height / (1 << hw.Level)

	return s.SoftThreshold(offsetX, offsetY, threshold)
}
