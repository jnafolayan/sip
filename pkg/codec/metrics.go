package codec

import (
	"image"
	"math"
)

type CompressionResult struct {
	Ratio float64
	PSNR  float64
}

func computeCompressionResult(img1, img2 image.Image) CompressionResult {
	return CompressionResult{
		PSNR: calcPSNR(img1, img2),
	}
}

func calcPSNR(img1, img2 image.Image) float64 {
	mse := calcMeanSquaredError(img1, img2)
	if mse == 0 {
		return math.Inf(1)
	}

	return 10 * math.Log10((255*255)/mse)
}

func calcMeanSquaredError(img1, img2 image.Image) float64 {
	bounds := img1.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var sum, mse float64

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Normalize the color values to [0, 255]
			rr1, gg1, bb1 := float64(r1>>8), float64(g1>>8), float64(b1>>8)
			rr2, gg2, bb2 := float64(r2>>8), float64(g2>>8), float64(b2>>8)

			// Calculate the squared error for each color channel
			sum += math.Pow(rr1-rr2, 2)
			sum += math.Pow(gg1-gg2, 2)
			sum += math.Pow(bb1-bb2, 2)
		}
	}

	mse = sum / float64(width*height*3)

	return mse
}

func computeCompressionResultBetweenImageData(a, b []uint8, width, height int) CompressionResult {
	return CompressionResult{
		PSNR: calcPSNRBetweenImageData(a, b, width, height),
	}
}

func calcPSNRBetweenImageData(a, b []uint8, width, height int) float64 {
	mse := calcMeanSquaredErrorBetweenImageData(a, b, width, height)
	if mse == 0 {
		return math.Inf(1)
	}

	return 10 * math.Log10((255*255)/mse)
}

func calcMeanSquaredErrorBetweenImageData(a, b []uint8, width, height int) float64 {
	var sum, mse float64
	var offset int
	var r1, g1, b1 float64
	var r2, g2, b2 float64

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset = (x + y*width) * 4
			r1 = float64(a[offset+0])
			g1 = float64(a[offset+1])
			b1 = float64(a[offset+2])
			r2 = float64(b[offset+0])
			g2 = float64(b[offset+1])
			b2 = float64(b[offset+2])

			// Calculate the squared error for each color channel
			sum += math.Pow(r1-r2, 2)
			sum += math.Pow(g1-g2, 2)
			sum += math.Pow(b1-b2, 2)
		}
	}

	mse = sum / float64(width*height*3)

	return mse
}
