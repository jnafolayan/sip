package imageutils

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"os"

	"github.com/jnafolayan/sip/pkg/signal"
)

type SignalCoeff = signal.SignalCoeff

func ReadImage(fpath string) (image.Image, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func SaveImage(fpath string, img image.Image) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, nil)
}

func Grayscale(img image.Image) [][]SignalCoeff {
	size := img.Bounds().Size()
	grayscale := make([][]SignalCoeff, size.Y)

	for i := 0; i < size.Y; i++ {
		grayscale[i] = make([]SignalCoeff, size.X)
		for j := 0; j < size.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			r1 := SignalCoeff(r >> 8)
			g1 := SignalCoeff(g >> 8)
			b1 := SignalCoeff(b >> 8)
			grayscale[i][j] = (r1 + g1 + b1) / 3.0
			// y, _, _ := color.RGBToYCbCr(r1, g1, b1)
			// grayscale[i][j] = SignalCoeff(y)
		}
	}

	return grayscale
}

func YCbCr(img image.Image) [][]color.YCbCr {
	size := img.Bounds().Size()
	pixelData := make([][]color.YCbCr, size.Y)

	for i := 0; i < size.Y; i++ {
		pixelData[i] = make([]color.YCbCr, size.X)
		for j := 0; j < size.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			r1 := uint8(r >> 8)
			g1 := uint8(g >> 8)
			b1 := uint8(b >> 8)
			y, cb, cr := color.RGBToYCbCr(r1, g1, b1)
			pixelData[i][j] = color.YCbCr{y, cb, cr}
		}
	}

	return pixelData
}

func ExtractYCbCrComponents(pixels [][]color.YCbCr) ([][]SignalCoeff, [][]SignalCoeff, [][]SignalCoeff) {
	Y := make([][]SignalCoeff, len(pixels))
	Cb := make([][]SignalCoeff, len(pixels))
	Cr := make([][]SignalCoeff, len(pixels))
	for i := 0; i < len(pixels); i++ {
		Y[i] = make([]SignalCoeff, len(pixels[i]))
		Cb[i] = make([]SignalCoeff, len(pixels[i]))
		Cr[i] = make([]SignalCoeff, len(pixels[i]))
		for j := 0; j < len(pixels[i]); j++ {
			Y[i][j] = SignalCoeff(pixels[i][j].Y)
			Cb[i][j] = SignalCoeff(pixels[i][j].Cb)
			Cr[i][j] = SignalCoeff(pixels[i][j].Cr)
		}
	}

	return Y, Cb, Cr
}
