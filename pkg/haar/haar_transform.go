package haar

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

func (hw *HaarWavelet) Transform(s signal.Signal2D) signal.Signal2D {
	width, height := getPowerOf2Size(s, hw.Level)
	result := s.Clone().Pad(width, height, signal.PadSymmetric)

	tempSignal := signal.New(width, height)

	w := width / 2
	h := height / 2

	for level := 0; level < hw.Level; level++ {
		// Transform rows
		for i := 0; i < h*2; i++ {
			// Calculate the averages and differences repeatedly
			if w == 0 {
				break
			}
			for j := 0; j < w; j++ {
				tempSignal[i][j] = (result[i][j*2] + result[i][j*2+1]) / 2.0
				tempSignal[i][j+w] = result[i][j*2] - tempSignal[i][j]
			}
			for j := 0; j < w*2; j++ {
				result[i][j] = tempSignal[i][j]
			}
		}

		// Transform columns
		for j := 0; j < w*2; j++ {
			// Calculate the averages and differences repeatedly
			if h == 0 {
				break
			}
			for i := 0; i < h; i++ {
				tempSignal[i][j] = (result[i*2][j] + result[i*2+1][j]) / 2.0
				tempSignal[i+h][j] = result[i*2][j] - tempSignal[i][j]
			}
			for i := 0; i < h*2; i++ {
				result[i][j] = tempSignal[i][j]
			}
		}

		w /= 2
		h /= 2
	}

	return result
}
