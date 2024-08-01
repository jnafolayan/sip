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

func ExtractYCbCrComponents(img image.Image) ([][]SignalCoeff, [][]SignalCoeff, [][]SignalCoeff) {
	size := img.Bounds().Size()

	Y := make([][]SignalCoeff, size.Y)
	Cb := make([][]SignalCoeff, size.Y)
	Cr := make([][]SignalCoeff, size.Y)

	for i := 0; i < size.Y; i++ {
		Y[i] = make([]SignalCoeff, size.X)
		Cb[i] = make([]SignalCoeff, size.X)
		Cr[i] = make([]SignalCoeff, size.X)
		for j := 0; j < size.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			r1 := uint8(r >> 8)
			g1 := uint8(g >> 8)
			b1 := uint8(b >> 8)
			y, cb, cr := color.RGBToYCbCr(r1, g1, b1)

			Y[i][j] = SignalCoeff(y)
			Cb[i][j] = SignalCoeff(cb)
			Cr[i][j] = SignalCoeff(cr)
		}
	}

	return Y, Cb, Cr
}

func ExtractYCbCrComponentsFromImageData(imageData []uint8, width, height int) ([][]SignalCoeff, [][]SignalCoeff, [][]SignalCoeff) {
	Y := make([][]SignalCoeff, height)
	Cb := make([][]SignalCoeff, height)
	Cr := make([][]SignalCoeff, height)

	var r, g, b uint8
	var offset int

	for y := 0; y < height; y++ {
		Y[y] = make([]SignalCoeff, width)
		Cb[y] = make([]SignalCoeff, width)
		Cr[y] = make([]SignalCoeff, width)
		for x := 0; x < width; x++ {
			offset = (x + y*width) * 4
			r = imageData[offset+0]
			g = imageData[offset+1]
			b = imageData[offset+2]
			yy, cb, cr := color.RGBToYCbCr(r, g, b)

			Y[y][x] = SignalCoeff(yy)
			Cb[y][x] = SignalCoeff(cb)
			Cr[y][x] = SignalCoeff(cr)
		}
	}

	return Y, Cb, Cr
}
