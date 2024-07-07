package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/haar"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

var compressFlags = &(struct {
	waveletType string
	outputFile  string
	level       int
}{})

var compressCmd = &cli.Command{
	Name: "compress",
	Init: func(cmd *cli.Command) {
		cmd.FlagSet = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
		cmd.FlagSet.StringVar(&compressFlags.waveletType, "wavelet", "haar", "wavelet type")
		cmd.FlagSet.IntVar(&compressFlags.level, "level", 1, "level of decomposition")
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

		grayscale := imageutils.Grayscale(img)

		var w wavelet.Wavelet
		switch compressFlags.waveletType {
		case "haar":
			w = &haar.HaarWavelet{Level: compressFlags.level}
		default:
			return fmt.Errorf("unrecognized wavelet: %s", compressFlags.waveletType)
		}

		transformed := w.Transform(grayscale)
		output := transformed.Image()

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		outFile := filepath.Join(wd, compressFlags.outputFile)
		err = imageutils.SaveImage(outFile, output)

		return err
	},
}
