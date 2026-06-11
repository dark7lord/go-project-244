// Package main implements a CLI utility to get the difference between two files like git diff
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code"
	"code/formatters"
)

const (
	binaryName  = "gendiff"
	ExitOK      = 0
	ExitGeneral = 1
	ExitUsage   = 64
	ExitDataErr = 65
)

func Run() error {
	cmd := &cli.Command{
		Name:                  binaryName,
		Usage:                 "Compares two configuration files and shows a difference.",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				DefaultText: "stylish",
				Value:       "stylish",
				Usage:       "output format",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if argsCount := cmd.NArg(); argsCount != 2 {
				return fmt.Errorf("expected 2 file paths, got %d", argsCount)
			}

			format := cmd.String("format")
			if !formatters.IsValidFormat(format) {
				return fmt.Errorf("unsupported format: %s", format)
			}

			filepathA := cmd.Args().Get(0)
			filepathB := cmd.Args().Get(1)

			diff, err := code.GenDiff(filepathA, filepathB, format)
			if err != nil {
				return err
			}

			fmt.Print(diff)

			return nil
		},
	}

	return cmd.Run(context.Background(), os.Args)
}

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
