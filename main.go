package main

import (
	"fmt"
	"os"

	"github.com/jnafolayan/sip/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Errorf("Error running command: %s", err)
		os.Exit(1)
	}
}
