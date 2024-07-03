package cmd

import (
	"github.com/jnafolayan/sip/internal/cli"
)

var rootCmd = &cli.Command{
	Name: "sip",
}

func Execute() error {
	return rootCmd.Execute()
}
