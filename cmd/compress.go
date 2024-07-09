package cmd

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/haar"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

var compressFlags = &(struct {
	waveletType string
	outputFile  string
	level       int
	threshold   int
}{})

var compressCmd = &cli.Command{
	Name: "compress",
	Init: func(cmd *cli.Command) {
		cmd.FlagSet = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
		cmd.FlagSet.StringVar(&compressFlags.waveletType, "wavelet", "haar", "wavelet type")
		cmd.FlagSet.IntVar(&compressFlags.level, "level", 1, "level of decomposition")
		cmd.FlagSet.IntVar(&compressFlags.threshold, "threshold", 0, "threshold")
		cmd.FlagSet.StringVar(&compressFlags.outputFile, "output", "", "output file")
	},
	Run: func(cmd *cli.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("compress: no image supplied")
		}

		img, err := imageutils.ReadImage(args[0])
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}

		// Parse flags
		cmd.FlagSet.Parse(args[1:])

		var w wavelet.Wavelet
		switch compressFlags.waveletType {
		case "haar":
			w = &haar.HaarWavelet{Level: compressFlags.level}
		default:
			return fmt.Errorf("unrecognized wavelet: %s", compressFlags.waveletType)
		}

		yCbCrPixels := imageutils.YCbCr(img)
		Y, Cb, Cr := imageutils.ExtractYCbCrComponents(yCbCrPixels)

		// Transform channels
		tY, tCb, tCr := transformYCbCr(w, Y, Cb, Cr)
		tWidth, tHeight := tY.Size()

		// Threshold channels
		offsetX := tWidth / (1 << w.GetDecompositionLevel())
		offsetY := tHeight / (1 << w.GetDecompositionLevel())
		threshold := compressFlags.threshold
		tY = tY.HardThreshold(offsetX, offsetY, threshold)
		tCb = tCb.HardThreshold(offsetX, offsetY, threshold)
		tCr = tCr.HardThreshold(offsetX, offsetY, threshold)

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		// transformedImage := createImageFromYCbCr(tY, tCb, tCr)
		outY, outCb, outCr := inverseTransformYCbCr(w, tY, tCb, tCr)

		// Remove artifacts caused by padding image
		width, height := len(Y[0]), len(Y)
		outY = outY.Slice(0, 0, width, height)
		outCb = outCb.Slice(0, 0, width, height)
		outCr = outCr.Slice(0, 0, width, height)

		output := createImageFromYCbCr(outY, outCb, outCr)
		originalImage := createImageFromYCbCr(Y, Cb, Cr)

		outFile := filepath.Join(wd, compressFlags.outputFile)
		err = imageutils.SaveImage(outFile, output)
		if err != nil {
			return err
		}

		psnr := calcPSNR(originalImage, output)
		fmt.Printf("Peak Signal-to-Noise ratio: %.2f\n", psnr)

		// fmt.Println(signal.Signal2D(tY).String(image.Rect(0, 0, 10, 10)))
		// fmt.Println(outY.Equal(Y))
		// fmt.Println(outCb.Equal(Cb))
		// fmt.Println(outCr.Equal(Cr))
		// fmt.Println(outY.String(image.Rect(0, 0, 10, 10)))

		return nil
	},
}

func transformYCbCr(
	w wavelet.Wavelet,
	Y, Cb, Cr signal.Signal2D,
) (signal.Signal2D, signal.Signal2D, signal.Signal2D) {
	YChan := make(chan signal.Signal2D, 1)
	CbChan := make(chan signal.Signal2D, 1)
	CrChan := make(chan signal.Signal2D, 1)

	go compressChannel(w, Y, YChan)
	go compressChannel(w, Cb, CbChan)
	go compressChannel(w, Cr, CrChan)

	tY := <-YChan
	tCb := <-CbChan
	tCr := <-CrChan

	return tY, tCb, tCr
}

func compressChannel(w wavelet.Wavelet, channel signal.Signal2D, out chan<- signal.Signal2D) {
	transformed := w.Transform(channel)
	out <- transformed
	close(out)
}

func inverseTransformYCbCr(
	w wavelet.Wavelet,
	tY, tCb, tCr signal.Signal2D,
) (signal.Signal2D, signal.Signal2D, signal.Signal2D) {
	YChan := make(chan signal.Signal2D, 1)
	CbChan := make(chan signal.Signal2D, 1)
	CrChan := make(chan signal.Signal2D, 1)

	go decompressChannel(w, tY, YChan)
	go decompressChannel(w, tCb, CbChan)
	go decompressChannel(w, tCr, CrChan)

	Y := <-YChan
	Cb := <-CbChan
	Cr := <-CrChan

	return Y, Cb, Cr
}

func decompressChannel(w wavelet.Wavelet, channel signal.Signal2D, out chan<- signal.Signal2D) {
	transformed := w.InverseTransform(channel)
	out <- transformed
	close(out)
}

func createImageFromYCbCr(Y, Cb, Cr signal.Signal2D) image.Image {
	width, height := Y.Size()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b := color.YCbCrToRGB(uint8(Y[i][j]), uint8(Cb[i][j]), uint8(Cr[i][j]))
			c := color.RGBA{r, g, b, 255}
			img.Set(j, i, c)
		}
	}

	return img
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

	var sum float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Normalize the color values to [0, 255]
			r1, g1, b1 = r1>>8, g1>>8, b1>>8
			r2, g2, b2 = r2>>8, g2>>8, b2>>8

			// Calculate the squared error for each color channel
			sum += math.Pow(float64(r1-r2), 2)
			sum += math.Pow(float64(g1-g2), 2)
			sum += math.Pow(float64(b1-b2), 2)
		}
	}

	return sum / float64(width*height*3)
}
