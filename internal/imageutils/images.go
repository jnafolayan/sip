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

func SaveImage(dest string, img image.Image) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 75})
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

func ConvertImageToImageData(img image.Image) []uint8 {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	result := make([]uint8, width*height*4)

	var r, g, b, a uint32
	var offset int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset = (x + y*width) * 4
			r, g, b, a = img.At(x, y).RGBA()
			result[offset+0] = uint8(r >> 8)
			result[offset+1] = uint8(g >> 8)
			result[offset+2] = uint8(b >> 8)
			result[offset+3] = uint8(a >> 8)
		}
	}

	return result
}

func ConvertImageDataToImage(imageData []uint8, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var r, g, b, a uint8
	var x, y int

	for i := 0; i < len(imageData); i += 4 {
		r = imageData[i+0]
		g = imageData[i+1]
		b = imageData[i+2]
		a = imageData[i+3]
		x = (i / 4) % width
		y = (i / 4) / width
		img.Set(x, y, color.RGBA{r, g, b, a})
	}

	return img
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

func GetImageChannelsFromImageData(imageData []uint8, width, height int) []signal.Signal2D {
	Y, Cb, Cr := ExtractYCbCrComponentsFromImageData(imageData, width, height)
	channels := []signal.Signal2D{Y, Cb, Cr}
	return channels
}

func GetImageChannels(img image.Image) []signal.Signal2D {
	Y, Cb, Cr := ExtractYCbCrComponents(img)
	channels := []signal.Signal2D{Y, Cb, Cr}
	return channels
}

func ReconstructImage(channels []signal.Signal2D, src image.Image) image.Image {
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

func ReconstructImageWithAlpha(channels []signal.Signal2D, alpha uint8) image.Image {
	width, height := channels[0].Size()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	Y, Cb, Cr := channels[0], channels[1], channels[2]

	var r, g, b uint8
	var yy, cb, cr float64
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

			r, g, b = color.YCbCrToRGB(uint8(yy), uint8(cb), uint8(cr))
			c = color.RGBA{r, g, b, alpha}
			img.Set(x, y, c)
		}
	}

	return img
}

func ReconstructImageData(channels []signal.Signal2D, original []uint8, out []uint8) []uint8 {
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
