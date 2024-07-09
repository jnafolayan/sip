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

func (hw *HaarWavelet) Transform(s signal.Signal2D) signal.Signal2D {
	width, height := getPowerOf2Size(s, hw.Level)
	result := s.Clone().Pad(width, height, signal.PadSymmetric)

	tempSignal := signal.New(width, height)

	// Transform rows
	for i := 0; i < height; i++ {
		// Calculate the averages and differences repeatedly
		level := 1
		for level <= hw.Level {
			w := width / (1 << level)
			if w == 0 {
				break
			}
			for j := 0; j < w; j++ {
				tempSignal[i][j] = (result[i][j*2] + result[i][j*2+1]) / 2.0
				tempSignal[i][j+w] = result[i][j*2] - tempSignal[i][j]
			}
			copy(result[i], tempSignal[i])
			level++
		}
	}

	// Transform columns
	for j := 0; j < width; j++ {
		// Calculate the averages and differences repeatedly
		level := 1
		for level <= hw.Level {
			h := height / (1 << level)
			if h == 0 {
				break
			}
			for i := 0; i < h; i++ {
				tempSignal[i][j] = (result[i*2][j] + result[i*2+1][j]) / 2.0
				tempSignal[i+h][j] = result[i*2][j] - tempSignal[i][j]
			}
			for i := 0; i < height; i++ {
				result[i][j] = tempSignal[i][j]
			}
			level++
		}
	}

	return result
}

func (hw *HaarWavelet) InverseTransform(s signal.Signal2D) signal.Signal2D {
	width, height := s.Size()
	result := s.Clone()

	tempSignal := s.Clone()

	for i := 0; i < height; i++ {
		level := hw.Level
		for level > 0 {
			w := width / (1 << level)
			for j := 0; j < w; j++ {
				tempSignal[i][j*2] = result[i][j] + result[i][j+w]
				tempSignal[i][j*2+1] = result[i][j] - result[i][j+w]
			}
			copy(result[i], tempSignal[i])
			level--
		}
	}

	for j := 0; j < width; j++ {
		level := hw.Level
		for level > 0 {
			h := height / (1 << level)
			for i := 0; i < h; i++ {
				tempSignal[i*2][j] = result[i][j] + result[i+h][j]
				tempSignal[i*2+1][j] = result[i][j] - result[i+h][j]
			}
			for i := 0; i < height; i++ {
				result[i][j] = tempSignal[i][j]
			}
			level--
		}
	}

	return result
}
