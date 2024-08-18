package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

var bulkCFlags = &(struct {
	source string
}{})

var bulkCompressCmd = &cli.Command{
	Name: "bulkc",
	Init: func(cmd *cli.Command) {
		cmd.FlagSet = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
		cmd.FlagSet.StringVar(&bulkCFlags.source, "src", "", "path to image")
	},
	Run: func(cmd *cli.Command, args []string) error {
		if len(args) == 0 {
			cmd.FlagSet.Usage()
			return fmt.Errorf("bulkc: no config supplied")
		}

		if err := cmd.FlagSet.Parse(args); err != nil {
			return err
		}

		src := bulkCFlags.source

		thresholds := []int{5, 15, 30}
		wavelets := []string{"haar", "cdf97"}

		rows := map[string]codec.CompressionResult{}

		os.RemoveAll("dump/")
		os.Mkdir("dump", 0777)

		for level := 1; level <= 6; level += 2 {
			for _, thresh := range thresholds {
				for _, wave := range wavelets {
					dest := fmt.Sprintf("dump/%dL_%dT_%sW.jpg", level, thresh, wave)
					fmt.Printf("Compressing %s\n", dest)
					result, err := codec.EncodeFileAsJPEG(src, dest, codec.CodecOptions{
						Wavelet:              wavelet.WaveletType(wave),
						ThresholdingFactor:   thresh,
						DecompositionLevel:   level,
						ThresholdingStrategy: "hard",
					})
					if err != nil {
						return err
					}
					rows[dest] = result
				}
			}
		}

		f, err := os.Create("dump/bulk.json")
		if err != nil {
			return err
		}

		return json.NewEncoder(f).Encode(rows)
	},
}
