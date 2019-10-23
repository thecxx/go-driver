package main

import (
	"os"

	"github.com/thecxx/go-driver/pkg/builder/cmd"
)

func main() {
	command := cmd.NewBuilderCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
