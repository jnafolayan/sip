package codec

import (
	"image"
)

type CompressionResult struct {
	Ratio float64
	PSNR  float64
	SSIM  float64
}

func computeCompressionResult(img1, img2 image.Image) CompressionResult {
	return CompressionResult{
		PSNR: calcPSNR(img1, img2),
		SSIM: calcSSIM(img1, img2),
	}
}

func computeCompressionResultBetweenImageData(a, b []uint8, width, height int) CompressionResult {
	return CompressionResult{
		PSNR: calcPSNRBetweenImageData(a, b, width, height),
		SSIM: calcSSIMBetweenImageData(a, b, width, height),
	}
}
