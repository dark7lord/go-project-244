// Package main implements a CLI utility for comparing two configuration files (JSON/YAML)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code"
)

const binaryName = "gendiff"

const (
	ExitOK    = 0
	ExitError = 1
)

func run(args []string) error {
	cmd := &cli.Command{
		Name:                  binaryName,
		Usage:                 "Compares two configuration files (JSON/YAML) and shows a difference (stylish, plain, json).",
		ArgsUsage:             "<old_file> <new_file>",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				DefaultText: "stylish",
				Value:       "stylish",
				Usage:       "output format (stylish, plain, json)",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if argsCount := cmd.NArg(); argsCount != 2 {
				return fmt.Errorf("expected 2 file paths, got %d", argsCount)
			}

			format := cmd.String("format")
			pathOldFile := cmd.Args().Get(0)
			pathNewFile := cmd.Args().Get(1)

			diff, err := code.GenDiff(pathOldFile, pathNewFile, format)
			if err != nil {
				return err
			}

			fmt.Print(diff)

			return nil
		},
	}

	return cmd.Run(context.Background(), args)
}

var osExit = os.Exit

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		osExit(ExitError)
		return
	}
	osExit(ExitOK)
}
