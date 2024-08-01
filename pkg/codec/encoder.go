package codec

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/cdf97"
	"github.com/jnafolayan/sip/pkg/ezw"
	"github.com/jnafolayan/sip/pkg/haar"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

func EncodeFileAsJPEG(source string, out string, opts CodecOptions) (CompressionResult, error) {
	img, err := imageutils.ReadImage(source)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("compress: %w", err)
	}

	result, err := EncodeAsJPEG(img, out, opts)
	if err == nil {
		srcStat, srcStatErr := os.Stat(source)
		outStat, outStatErr := os.Stat(out)
		if srcStatErr != nil || outStatErr != nil {
			result.Ratio = -1
			return result, nil
		}

		result.Ratio = float64(srcStat.Size()) / float64(outStat.Size())
	}

	return result, err
}

func EncodeAsJPEG(img image.Image, out string, opts CodecOptions) (CompressionResult, error) {
	compressed, result := Encode(img, opts)

	err := imageutils.SaveImage(out, compressed)
	if err != nil {
		return result, err
	}

	return result, err
}

func Encode(img image.Image, opts CodecOptions) (image.Image, CompressionResult) {
	w, _ := getWaveletFamily(opts.Wavelet, opts)
	imageChannels := getImageChannels(img)

	channels := imageChannels
	channels = transformChannels(w, opts.ThresholdingFactor, channels)

	// Inverse transform
	channels = inverseTransformChannels(w, channels)

	// Remove padding that might be added during the transform stage
	imgBounds := img.Bounds().Size()
	originalWidth, originalHeight := imgBounds.X, imgBounds.Y
	channels = trimFatChannels(channels, originalWidth, originalHeight)

	reconstructed := reconstructYCbCrImage(channels)

	// Reconstruct the original image from the unprocessed (original) channels. Doing this
	// to overcome precision that might have been lost due to conversion from RGB to YCbCr.
	originalImage := reconstructYCbCrImage(imageChannels)
	result := computeCompressionResult(originalImage, reconstructed)

	return reconstructed, result
}

func transformChannels(w wavelet.Wavelet, threshold int, channels []signal.Signal2D) []signal.Signal2D {
	transformedChannels := make([]signal.Signal2D, len(channels))
	for i, c := range channels {
		// Apply wavelet transform on the channels
		transformed := w.Transform(c)
		// Hard thresholding
		coeffs := w.HardThreshold(transformed, threshold)
		transformedChannels[i] = coeffs
	}

	return transformedChannels
}

func inverseTransformChannels(w wavelet.Wavelet, channels []signal.Signal2D) []signal.Signal2D {
	transformedChannels := make([]signal.Signal2D, len(channels))
	for i, c := range channels {
		transformed := w.InverseTransform(c)
		transformedChannels[i] = transformed
	}

	return transformedChannels
}

func trimFatChannels(channels []signal.Signal2D, w, h int) []signal.Signal2D {
	trimmedChannels := make([]signal.Signal2D, len(channels))
	for i, c := range channels {
		trimmedChannels[i] = c.Slice(0, 0, w, h)
	}
	return trimmedChannels
}

func reconstructImage(channels []signal.Signal2D, joinChannelAtPos func(int, int) color.Color) image.Image {
	width, height := channels[0].Size()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, joinChannelAtPos(x, y))
		}
	}

	return img
}

func reconstructYCbCrImage(channels []signal.Signal2D) image.Image {
	Y, Cb, Cr := channels[0], channels[1], channels[2]
	return reconstructImage(channels, func(x, y int) color.Color {
		r, g, b := color.YCbCrToRGB(uint8(Y[y][x]), uint8(Cb[y][x]), uint8(Cr[y][x]))
		return color.RGBA{r, g, b, 255}
	})
}

func getWaveletFamily(wType wavelet.WaveletType, opts CodecOptions) (wavelet.Wavelet, error) {
	// Get wavelet family
	var w wavelet.Wavelet
	switch wType {
	case wavelet.WaveletHaar:
		w = &haar.HaarWavelet{Level: opts.DecompositionLevel}
	case wavelet.WaveletCDF97:
		w = &cdf97.CDF97Wavelet{Level: opts.DecompositionLevel}
	default:
		return nil, fmt.Errorf("unrecognized wavelet: %s", wType)
	}
	return w, nil
}

// createEncoders creates and initialized EZW encoders for each
// channel in an image.
// Currently there are 3 channels: Y, Cb and Cr, adopting the YCbCr
// color model.
func createEncoders(channels []signal.Signal2D, opts CodecOptions) []*ezw.Encoder {
	encoders := make([]*ezw.Encoder, len(channels))

	for i := range channels {
		e := ezw.NewEncoder()
		e.Init(channels[i], opts.DecompositionLevel)
		encoders[i] = e
	}

	return encoders
}

func getImageChannels(img image.Image) []signal.Signal2D {
	yCbCrPixels := imageutils.YCbCr(img)
	Y, Cb, Cr := imageutils.ExtractYCbCrComponents(yCbCrPixels)

	channels := []signal.Signal2D{Y, Cb, Cr}
	return channels
}
