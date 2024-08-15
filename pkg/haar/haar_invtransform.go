package haar

import "github.com/jnafolayan/sip/pkg/signal"

func (hw *HaarWavelet) InverseTransform(s signal.Signal2D) signal.Signal2D {
	width, height := s.Size()
	result := s.Clone()

	tempSignal := s.Clone()

	w := width / (1 << hw.Level)
	h := height / (1 << hw.Level)

	for level := 0; level < hw.Level; level++ {
		for i := 0; i < h*2; i++ {
			for j := 0; j < w; j++ {
				tempSignal[i][j*2] = result[i][j] + result[i][j+w]
				tempSignal[i][j*2+1] = result[i][j] - result[i][j+w]
			}
			for j := 0; j < w*2; j++ {
				result[i][j] = tempSignal[i][j]
			}
		}

		for j := 0; j < w*2; j++ {
			for i := 0; i < h; i++ {
				tempSignal[i*2][j] = result[i][j] + result[i+h][j]
				tempSignal[i*2+1][j] = result[i][j] - result[i+h][j]
			}
			for i := 0; i < h*2; i++ {
				result[i][j] = tempSignal[i][j]
			}
		}

		w *= 2
		h *= 2
	}

	return result
}
