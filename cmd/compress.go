package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

var compressFlags = &(struct {
	waveletType string
	outputFile  string
	level       int
	threshold   int
	T           string
}{})

var compressCmd = &cli.Command{
	Name: "compress",
	Init: func(cmd *cli.Command) {
		cmd.FlagSet = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
		cmd.FlagSet.StringVar(&compressFlags.waveletType, "wavelet", "haar", "wavelet type")
		cmd.FlagSet.IntVar(&compressFlags.level, "level", 1, "level of decomposition")
		cmd.FlagSet.IntVar(&compressFlags.threshold, "threshold", 10, "threshold")
		cmd.FlagSet.StringVar(&compressFlags.T, "T", "hard", "thresholding strategy (soft or hard)")
		cmd.FlagSet.StringVar(&compressFlags.outputFile, "output", "", "output file")
	},
	Run: func(cmd *cli.Command, args []string) error {
		if len(args) == 0 {
			cmd.FlagSet.Usage()
			return fmt.Errorf("compress: no image supplied")
		}

		sourceFile := args[0]

		// Parse flags
		if err := cmd.FlagSet.Parse(args[1:]); err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		destFile := filepath.Join(wd, compressFlags.outputFile)

		result, err := codec.EncodeFileAsJPEG(sourceFile, destFile, codec.CodecOptions{
			Wavelet:              wavelet.WaveletType(compressFlags.waveletType),
			ThresholdingFactor:   compressFlags.threshold,
			DecompositionLevel:   compressFlags.level,
			ThresholdingStrategy: compressFlags.T,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Took %fs to compress.\n", result.Time)
		fmt.Println("PSNR:", result.PSNR)
		fmt.Println("SSIM:", result.SSIM)
		fmt.Println("Ratio:", result.Ratio)

		return nil
	},
}
