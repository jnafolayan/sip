package imageutils

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

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

func Grayscale(img image.Image) [][]float32 {
	size := img.Bounds().Size()
	grayscale := make([][]float32, size.Y)

	for i := 0; i < size.Y; i++ {
		grayscale[i] = make([]float32, size.X)
		for j := 0; j < size.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			r1 := float32(r >> 8)
			g1 := float32(g >> 8)
			b1 := float32(b >> 8)
			grayscale[i][j] = (r1 + g1 + b1) / 3.0
			// y, _, _ := color.RGBToYCbCr(r1, g1, b1)
			// grayscale[i][j] = float32(y)
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

func ExtractYCbCrComponents(pixels [][]color.YCbCr) ([][]float32, [][]float32, [][]float32) {
	Y := make([][]float32, len(pixels))
	Cb := make([][]float32, len(pixels))
	Cr := make([][]float32, len(pixels))
	for i := 0; i < len(pixels); i++ {
		Y[i] = make([]float32, len(pixels[i]))
		Cb[i] = make([]float32, len(pixels[i]))
		Cr[i] = make([]float32, len(pixels[i]))
		for j := 0; j < len(pixels[i]); j++ {
			Y[i][j] = float32(pixels[i][j].Y)
			Cb[i][j] = float32(pixels[i][j].Cb)
			Cr[i][j] = float32(pixels[i][j].Cr)
		}
	}

	return Y, Cb, Cr
}
