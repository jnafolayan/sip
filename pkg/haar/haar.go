package haar

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

type HaarWavelet struct {
	Level int
}

func (hw *HaarWavelet) GetDecompositionLevel() int {
	return hw.Level
}

func getPowerOf2Size(s signal.Signal2D, level int) (int, int) {
	N, M := s.Size()
	if N != (N>>level)<<level {
		N = (N>>level + 1) << level
	}
	if M != (M>>level)<<level {
		M = (M>>level + 1) << level
	}
	return N, M
}
