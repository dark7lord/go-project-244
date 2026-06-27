// Package main provides the gendiff command-line entrypoint.
package main

import (
	"os"

	"code/internal/cli"
)

func main() {
	exitCode := cli.Execute(os.Args, os.Stderr)
	os.Exit(exitCode)
}
