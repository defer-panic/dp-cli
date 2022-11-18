package main

import (
	"os"

	"github.com/defer-panic/dp-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

