package imageutils

import (
	"image"
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
