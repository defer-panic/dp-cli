package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/defer-panic/dp-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

