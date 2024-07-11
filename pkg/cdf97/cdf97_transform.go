package cdf97

import "github.com/jnafolayan/sip/pkg/signal"

var analysisLowPassFilter = []float32{
	0.026748757411, -0.016864118443, -0.078223266529,
	0.266864118443, 0.602949018236, 0.266864118443,
	-0.078223266529, -0.016864118443, 0.026748757411,
}

var analysisHighPassFilter = []float32{
	0, 0.091271763114, -0.057543526229,
	-0.591271763114, 1.11508705, -0.591271763114,
	-0.057543526229, 0.091271763114, 0,
}

func convolve(s, filter []float32) []float32 {
	result := make([]float32, len(s)+len(filter)-1)
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(filter); j++ {
			result[i+j] += s[i] * filter[j]
		}
	}
	return result
}

func downsample(s []float32) []float32 {
	result := make([]float32, (len(s)+1)/2)
	for i := 0; i < len(s); i += 2 {
		result[i/2] = s[i]
	}
	return result
}

func downsampleCols(s signal.Signal2D) signal.Signal2D {
	width, height := s.Size()
	result := signal.New(width, (height+1)/2)
	colData := make([]float32, height)

	for x := 0; x < width; x++ {
		// Pull out column
		for y := 0; y < height; y++ {
			colData[y] = s[y][x]
		}
		downsampled := downsample(colData)
		// Patch the signal with the convolved column
		for y := 0; y < len(downsampled); y++ {
			result[y][x] = downsampled[y]
		}
	}

	return result
}

func applyFilterRows(s signal.Signal2D, filter []float32) signal.Signal2D {
	width, height := s.Size()
	result := signal.New(width+len(filter)-1, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result[y] = convolve(s[y], filter)
		}
	}
	return result
}

func applyFilterCols(s signal.Signal2D, filter []float32) signal.Signal2D {
	width, height := s.Size()
	result := signal.New(width, height+len(filter)-1)

	colData := make([]float32, height)
	for x := 0; x < width; x++ {
		// Pull out column
		for y := 0; y < height; y++ {
			colData[y] = s[y][x]
		}
		filtered := convolve(colData, filter)
		// Patch the signal with the convolved column
		for y := 0; y < len(filtered); y++ {
			result[y][x] = filtered[y]
		}
	}

	return result
}

func (cdf *CDF97Wavelet) Transform(s signal.Signal2D) signal.Signal2D {
	ll, _, _, _ := cdf.Decompose(s, cdf.Level)
	return ll
}

func (cdf *CDF97Wavelet) Decompose(
	s signal.Signal2D, level int,
) (signal.Signal2D, signal.Signal2D, signal.Signal2D, signal.Signal2D) {
	// Apply filters to rows
	rowsLowPass := applyFilterRows(s, analysisLowPassFilter)
	rowsHighPass := applyFilterRows(s, analysisHighPassFilter)

	// Downsample rows
	for i := range rowsLowPass {
		rowsLowPass[i] = downsample(rowsLowPass[i])
		rowsHighPass[i] = downsample(rowsHighPass[i])
	}

	// Apply filters to columns
	ll := applyFilterCols(rowsLowPass, analysisLowPassFilter)
	lh := applyFilterCols(rowsLowPass, analysisHighPassFilter)
	hl := applyFilterCols(rowsHighPass, analysisLowPassFilter)
	hh := applyFilterCols(rowsHighPass, analysisHighPassFilter)

	// Downsample columns
	ll = downsampleCols(ll)
	lh = downsampleCols(lh)
	hl = downsampleCols(hl)
	hh = downsampleCols(hh)

	if level > 1 {
		ll, _, _, _ = cdf.Decompose(ll, level-1)
	}

	return ll, lh, hl, hh
}
