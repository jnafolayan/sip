package cmd

import (
	"fmt"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/haar"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type compressFlags struct {
	waveletType string
}

var compressCmd = &cli.Command{
	Name: "compress",
	Run: func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("compress: no image supplied")
		}

		img, err := imageutils.ReadImage(args[0])
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}

		grayscale := imageutils.Grayscale(img)

		var w wavelet.Wavelet
		w = &haar.HaarWavelet{}

		return nil
	},
}
