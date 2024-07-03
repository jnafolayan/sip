package cmd

import (
	"fmt"

	"github.com/jnafolayan/sip/internal/cli"
	"github.com/jnafolayan/sip/internal/imageutils"
)

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

		fmt.Println(img)

		return nil
	},
}
