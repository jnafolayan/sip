package codec

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"
	"time"

	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

func EncodeFileAsJPEG(src string, dest string, opts CodecOptions) (CompressionResult, error) {
	img, err := imageutils.ReadImage(src)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("compress: %w", err)
	}

	result, err := EncodeAsJPEG(img, dest, opts)
	if err == nil {
		srcStat, srcStatErr := os.Stat(src)
		outStat, outStatErr := os.Stat(dest)
		if srcStatErr != nil || outStatErr != nil {
			result.Ratio = -1
			return result, nil
		}

		result.Ratio = float64(srcStat.Size()) / float64(outStat.Size())
	}

	return result, err
}

func EncodeAsJPEG(img image.Image, dest string, opts CodecOptions) (CompressionResult, error) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	// Speed up computations by <66% using a 1-D RGBA slice over an image
	imageData := imageutils.ConvertImageToImageData(img)

	start := time.Now()
	compressedImageData, result := EncodeImageData(imageData, width, height, opts)
	result.Time = time.Since(start).Seconds()

	compressed := imageutils.ConvertImageDataToImage(compressedImageData, width, height)
	err := imageutils.SaveImage(dest, compressed)
	if err != nil {
		return result, err
	}

	return result, err
}

func Encode(img image.Image, opts CodecOptions) (image.Image, CompressionResult) {
	w, _ := GetWaveletFamily(opts)
	imageChannels := imageutils.GetImageChannels(img)

	channels := imageChannels
	channels = transformChannels(w, opts, channels)

	// Inverse transform
	channels = inverseTransformChannels(w, channels)

	// Remove padding that might be added during the transform stage
	imgBounds := img.Bounds().Size()
	originalWidth, originalHeight := imgBounds.X, imgBounds.Y
	channels = trimFatChannels(channels, originalWidth, originalHeight)

	reconstructed := reconstructImage(channels, img)

	// Reconstruct the original image from the unprocessed (original) channels. Doing this
	// to overcome precision that might have been lost due to conversion from RGB to YCbCr.
	originalImage := reconstructImage(imageChannels, img)
	result := computeCompressionResult(originalImage, reconstructed)

	return reconstructed, result
}

func EncodeImageData(imageData []uint8, width, height int, opts CodecOptions) ([]uint8, CompressionResult) {
	w, _ := GetWaveletFamily(opts)
	imageChannels := imageutils.GetImageChannelsFromImageData(imageData, width, height)

	channels := imageChannels
	channels = transformChannels(w, opts, channels)
	// debug.DrawSignal2D(imageChannels[0], image.Rect(50, 150, 60, 160), "original_Y.jpg")

	// Inverse transform
	channels = inverseTransformChannels(w, channels)
	// debug.DrawSignal2D(channels[0], image.Rect(50, 150, 60, 160), "transformed_Y.jpg")

	// Remove padding that might be added during the transform stage
	channels = trimFatChannels(channels, width, height)

	reconstructed := make([]uint8, width*height*4)
	reconstructed = reconstructImageData(channels, imageData, reconstructed)

	// Reconstruct the original image from the unprocessed (original) channels. Doing this
	// to overcome precision that might have been lost due to conversion from RGB to YCbCr.
	originalImageData := make([]uint8, width*height*4)
	originalImageData = reconstructImageData(imageChannels, imageData, originalImageData)
	result := computeCompressionResultBetweenImageData(originalImageData, reconstructed, width, height)

	return reconstructed, result
}

func transformChannels(w wavelet.Wavelet, opts CodecOptions, channels []signal.Signal2D) []signal.Signal2D {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	transformedChannels := make([]signal.Signal2D, len(channels))
	// transformedChannels := channels
	for i, c := range channels {
		go func() {
			// Apply wavelet transform on the channels
			transformed := w.Transform(c)
			// Hard thresholding
			if opts.ThresholdingStrategy == "soft" {
				transformed = w.SoftThreshold(transformed, opts.ThresholdingFactor)
			} else if opts.ThresholdingStrategy == "hard" {
				transformed = w.HardThreshold(transformed, opts.ThresholdingFactor)
			}
			transformedChannels[i] = transformed
			wg.Done()
		}()
	}

	wg.Wait()

	return transformedChannels
}

func inverseTransformChannels(w wavelet.Wavelet, channels []signal.Signal2D) []signal.Signal2D {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	// transformedChannels := make([]signal.Signal2D, len(channels))
	transformedChannels := channels
	for i, c := range channels {
		go func() {
			transformed := w.InverseTransform(c)
			transformedChannels[i] = transformed
			wg.Done()
		}()
	}

	wg.Wait()

	return transformedChannels
}

func trimFatChannels(channels []signal.Signal2D, w, h int) []signal.Signal2D {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	// trimmedChannels := make([]signal.Signal2D, len(channels))
	trimmedChannels := channels
	for i, c := range channels {
		go func() {
			trimmedChannels[i] = c.Slice(0, 0, w, h)
			wg.Done()
		}()
	}

	wg.Wait()

	return trimmedChannels
}

func reconstructImage(channels []signal.Signal2D, src image.Image) image.Image {
	width, height := channels[0].Size()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	Y, Cb, Cr := channels[0], channels[1], channels[2]

	var r, g, b uint8
	var yy, cb, cr float64
	var alpha uint32
	var c color.RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			yy, cb, cr = Y[y][x], Cb[y][x], Cr[y][x]

			// Clamp the values between 0 ... 255
			if yy < 0 {
				yy = 0
			} else if yy > 255 {
				yy = 255
			}
			if cb < 0 {
				cb = 0
			} else if cb > 255 {
				cb = 255
			}
			if cr < 0 {
				cr = 0
			} else if cr > 255 {
				cr = 255
			}

			_, _, _, alpha = src.At(x, y).RGBA()
			r, g, b = color.YCbCrToRGB(uint8(yy), uint8(cb), uint8(cr))
			c = color.RGBA{r, g, b, uint8(alpha >> 8)}
			img.Set(x, y, c)
		}
	}

	return img
}

func reconstructImageData(channels []signal.Signal2D, original []uint8, out []uint8) []uint8 {
	width, height := channels[0].Size()
	Y, Cb, Cr := channels[0], channels[1], channels[2]
	var r, g, b uint8

	var offset int
	var yy, cb, cr float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset = (x + y*width) * 4
			yy, cb, cr = Y[y][x], Cb[y][x], Cr[y][x]

			// Clamp the values between 0 ... 255
			if yy < 0 {
				yy = 0
			} else if yy > 255 {
				yy = 255
			}
			if cb < 0 {
				cb = 0
			} else if cb > 255 {
				cb = 255
			}
			if cr < 0 {
				cr = 0
			} else if cr > 255 {
				cr = 255
			}

			r, g, b = color.YCbCrToRGB(uint8(yy), uint8(cb), uint8(cr))
			out[offset+0] = r
			out[offset+1] = g
			out[offset+2] = b
			out[offset+3] = original[offset+3]
		}
	}

	return out
}
