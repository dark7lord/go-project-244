// Package main implements a CLI utility to get the difference between two files like git diff
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code/internal/gendiff"
)

const binaryName = "gendiff"

func run(args []string) error {
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
			pathOldFile := cmd.Args().Get(0)
			pathNewFile := cmd.Args().Get(1)

			diff, err := gendiff.GenDiff(pathOldFile, pathNewFile, format)
			if err != nil {
				return err
			}

			fmt.Print(diff)

			return nil
		},
	}

	return cmd.Run(context.Background(), args)
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
