package main

import (
	"os"

	"github.com/drlau/akashi/cmd"
)

func main() {
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
