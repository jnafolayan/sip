package codec

import (
	"image"

	"github.com/jnafolayan/sip/internal/imageutils"
)

func calcSSIM(img1, img2 image.Image) float64 {
	data1 := imageutils.ConvertImageToImageData(img1)
	data2 := imageutils.ConvertImageToImageData(img2)

	return calcMeanSquaredErrorBetweenImageData(data1, data2, img1.Bounds().Dx(), img1.Bounds().Dy())
}

func calcSSIMBetweenImageData(a, b []uint8, width, height int) float64 {
	const (
		K1 = 0.01
		K2 = 0.03
		L  = 255 // bit depth
	)

	// Constants for SSIM calculation
	C1 := (K1 * L) * (K1 * L)
	C2 := (K2 * L) * (K2 * L)

	// Means
	r1Mean, g1Mean, b1Mean := pixelsMean(a, width, height)
	r2Mean, g2Mean, b2Mean := pixelsMean(b, width, height)

	// Variances and covariance
	r1Variance, g1Variance, b1Variance := pixelsVariance(a, width, height, r1Mean, g1Mean, b1Mean)
	r2Variance, g2Variance, b2Variance := pixelsVariance(b, width, height, r2Mean, g2Mean, b2Mean)
	rCov, gCov, bCov := pixelsCovariance(
		a, r1Mean, g1Mean, b1Mean,
		b, r2Mean, g2Mean, b2Mean,
		width, height,
	)

	// Calculate SSIM for each channel
	rSSIM := calcSSIMForChannel(C1, C2, r1Mean, r2Mean, r1Variance, r2Variance, rCov)
	gSSIM := calcSSIMForChannel(C1, C2, g1Mean, g2Mean, g1Variance, g2Variance, gCov)
	bSSIM := calcSSIMForChannel(C1, C2, b1Mean, b2Mean, b1Variance, b2Variance, bCov)

	// Average the SSIM values of the three channels to get the overall SSIM
	overallSSIM := (rSSIM + gSSIM + bSSIM) / 3.0

	return overallSSIM
}

func calcSSIMForChannel(C1, C2, mean1, mean2, variance1, variance2, covariance float64) float64 {
	// SSIM formula
	numerator := (2*mean1*mean2 + C1) * (2*covariance + C2)
	denominator := (mean1*mean1 + mean2*mean2 + C1) * (variance1 + variance2 + C2)

	return numerator / denominator
}

// Function to calculate the mean of pixel values in a channel
func pixelsMean(data []uint8, w, h int) (r, g, b float64) {
	numPixels := float64(w * h)

	for i := 0; i < len(data); i += 4 {
		r += float64(data[i+0])
		g += float64(data[i+1])
		b += float64(data[i+2])
	}

	r /= numPixels
	g /= numPixels
	b /= numPixels

	return r, g, b
}

// Function to calculate the variance of pixel values in a channel
func pixelsVariance(data []uint8, w, h int, rMean, gMean, bMean float64) (r, g, b float64) {
	var rDiff, gDiff, bDiff float64
	numPixels := float64(w * h)

	for i := 0; i < len(data); i += 4 {
		rDiff = float64(data[i+0]) - rMean
		gDiff = float64(data[i+1]) - gMean
		bDiff = float64(data[i+2]) - bMean
		r += rDiff * rDiff
		g += gDiff * gDiff
		b += bDiff * bDiff
	}

	r /= numPixels
	g /= numPixels
	b /= numPixels

	return r, g, b
}

// Function to calculate the covariance between two channels
func pixelsCovariance(
	channel1 []uint8, r1Mean, g1Mean, b1Mean float64,
	channel2 []uint8, r2Mean, g2Mean, b2Mean float64,
	w, h int,
) (r, g, b float64) {
	var r1Diff, g1Diff, b1Diff float64
	var r2Diff, g2Diff, b2Diff float64
	numPixels := float64(w * h)
	for i := 0; i < len(channel1); i += 4 {
		r1Diff = float64(channel1[i+0]) - r1Mean
		g1Diff = float64(channel1[i+1]) - g1Mean
		b1Diff = float64(channel1[i+2]) - b1Mean

		r2Diff = float64(channel2[i+0]) - r2Mean
		g2Diff = float64(channel2[i+1]) - g2Mean
		b2Diff = float64(channel2[i+2]) - b2Mean

		r += r1Diff * r2Diff
		g += g1Diff * g2Diff
		b += b1Diff * b2Diff
	}

	r /= numPixels
	g /= numPixels
	b /= numPixels

	return r, g, b
}
