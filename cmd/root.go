package cmd

import (
	"flag"

	"github.com/jnafolayan/sip/internal/cli"
)

var rootCmd = &cli.Command{
	Name: "sip",
}

func init() {
	rootCmd.RegisterCmd(compressCmd)
	rootCmd.RegisterCmd(bulkCompressCmd)
}

func Execute() error {
	flag.Parse()
	return rootCmd.Execute(flag.Args())
}
