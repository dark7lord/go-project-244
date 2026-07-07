// Package main provides the gendiff command-line entrypoint.
package main

import (
	"context"
	"fmt"
	"os"

	"code/internal/cli"
)

func main() {
	cmd := cli.NewCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitError)
	}
}
