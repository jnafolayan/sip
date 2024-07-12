package cdf97

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

// getNearestMultipleOf2Size return the nearest largest multiple of 2 for each
// signal dimension.
func getNearestMultipleOf2Size(s signal.Signal2D, level int) (int, int) {
	N, M := s.Size()
	if N != (N>>level)<<level {
		N = (N>>level + 1) << level
	}
	if M != (M>>level)<<level {
		M = (M>>level + 1) << level
	}
	if N > M {
		M = N
	} else {
		N = M
	}
	return N, M
}

func (cdf *CDF97Wavelet) Transform(s signal.Signal2D) signal.Signal2D {
	width, height := getNearestMultipleOf2Size(s, cdf.Level)
	result := s.Clone()
	result = result.Pad(width, height, signal.PadSymmetric)

	for level := 0; level < cdf.Level; level++ {
		// Cols
		result = cdf.Decompose(result, width, height)
		result = transpose(result)
		// Rows
		result = cdf.Decompose(result, height, width)
		result = transpose(result)

		width /= 2
		height /= 2
	}

	return result
}

func (cdf *CDF97Wavelet) Decompose(s signal.Signal2D, width, height int) signal.Signal2D {
	// 9/7 Coefficients:
	var (
		a1 float32 = -1.586134342
		a2 float32 = -0.05298011854
		a3 float32 = 0.8829110762
		a4 float32 = 0.4435068522

		// Scale coeff:
		k1 float32 = 0.81289306611596146 // 1/1.230174104914
		k2 float32 = 0.61508705245700002 // 1.230174104914/2
	)

	// if width < 2 || height < 2 {
	// 	return s
	// }

	h1 := height - 1
	if h1 < 0 {
		h1 += height
	}
	h2 := height - 2
	if h2 < 0 {
		h2 += height
	}

	// Do 1D transform on all columns
	for x := 0; x < width; x++ {
		// Predict 1.
		for y := 1; y < height-1; y += 2 {
			s[y][x] += a1 * (s[y-1][x] + s[y+1][x])
		}
		s[h1][x] += 2 * a1 * s[h2][x] // Symmetric extension

		// Update 1.
		for y := 2; y < height; y += 2 {
			s[y][x] += a2 * (s[y-1][x] + s[y+1][x])
		}
		s[0][x] += 2 * a2 * s[1][x] // Symmetric extension

		// Predict 2.
		for y := 1; y < height-1; y += 2 {
			s[y][x] += a3 * (s[y-1][x] + s[y+1][x])
		}
		s[h1][x] += 2 * a3 * s[h2][x]

		// Update 2.
		for y := 2; y < height; y += 2 {
			s[y][x] += a4 * (s[y-1][x] + s[y+1][x])
		}
		s[0][x] += 2 * a4 * s[1][x]
	}

	// De-interleave
	tempBank := signal.New(height, width)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// k1 and k2 scale the vals
			// Simultaneously transpose the matrix when deinterleaving
			if y%2 == 0 {
				tempBank[x][y/2] = k1 * s[y][x]
			} else {
				tempBank[x][y/2+height/2] = k2 * s[y][x]
			}
		}
	}

	// Write tempBank to s
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			s[y][x] = tempBank[y][x]
		}
	}

	return s
}
