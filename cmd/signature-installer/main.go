package main

import (
	"fmt"
	"os"

	"outlook-signature/pkg/cli"
)

func main() {
	app := cli.App()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
