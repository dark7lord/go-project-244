package main

import (
	"os"

	"code/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args, os.Stderr))
}
