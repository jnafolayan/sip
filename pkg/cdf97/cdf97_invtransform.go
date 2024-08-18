package cdf97

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

func (cdf *CDF97Wavelet) InverseTransform(s signal.Signal2D) signal.Signal2D {
	width, height := s.Size()
	result := s.Clone()

	// fmt.Println(result.String(result.Bounds()) + "\n")
	// Find starting size of m:
	width >>= cdf.Level - 1
	height >>= cdf.Level - 1

	for level := 0; level < cdf.Level; level++ {
		// Rows
		// transposeInPlace(result)
		result = cdf.Reconstruct(result, width, height)
		// Cols
		// transposeInPlace(result)
		result = cdf.Reconstruct(result, width, height)

		width *= 2
		height *= 2
	}

	// Clamp the values between 0 ... 255
	width, height = s.Size()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if result[y][x] < 0 {
				result[y][x] = 0
			} else if result[y][x] > 255 {
				result[y][x] = 255
			}
		}
	}

	return result
}

func (cdf *CDF97Wavelet) Reconstruct(s signal.Signal2D, width, height int) signal.Signal2D {
	// 9/7 inverse coefficients:
	var (
		a1 signal.SignalCoeff = 1.586134342
		a2 signal.SignalCoeff = 0.05298011854
		a3 signal.SignalCoeff = -0.8829110762
		a4 signal.SignalCoeff = -0.4435068522

		// Inverse scale coeffs:
		k1 signal.SignalCoeff = 1.230174104914
		k2 signal.SignalCoeff = 1.6257861322319229
	)

	if width < 2 || height < 2 {
		return s
	}

	h1 := height - 1
	if h1 < 0 {
		h1 += height
	}
	h2 := height - 2
	if h2 < 0 {
		h2 += height
	}

	// Interleave
	tempBank := signal.New(width, height)
	for x := 0; x < width/2; x++ {
		for y := 0; y < height; y++ {
			// k1 and k2 scale the vals
			// Simultaneously transpose the matrix when interleaving
			tempBank[x*2][y] = k1 * s[y][x]
			tempBank[x*2+1][y] = k2 * s[y][x+width/2]
		}
	}

	// Write tempBank to s
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			s[y][x] = tempBank[y][x]
		}
	}

	// Do the 1D transform on all cols
	for x := 0; x < width; x++ {
		// Perform the inverse 1D transform.
		// Inverse update 2.
		for y := 2; y < height; y += 2 {
			s[y][x] += a4 * (s[y-1][x] + s[y+1][x])
		}
		s[0][x] += 2 * a4 * s[1][x]

		// Inverse predict 2.
		for y := 1; y < height-1; y += 2 {
			s[y][x] += a3 * (s[y-1][x] + s[y+1][x])
		}
		s[h1][x] += 2 * a3 * s[h2][x]

		// Inverse update 1.
		for y := 2; y < height; y += 2 {
			s[y][x] += a2 * (s[y-1][x] + s[y+1][x])
		}
		s[0][x] += 2 * a2 * s[1][x] // Symmetric extension

		// Inverse predict 1.
		for y := 1; y < height-1; y += 2 {
			s[y][x] += a1 * (s[y-1][x] + s[y+1][x])
		}
		s[h1][x] += 2 * a1 * s[h2][x] // Symmetric extension
	}

	return s
}
