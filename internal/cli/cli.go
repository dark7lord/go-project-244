package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/urfave/cli/v3"

	"code"
)

const BinaryName = "gendiff"
const (
	ExitOK    = 0
	ExitError = 1
)

func Run(args []string) error {
	cmd := &cli.Command{
		Name:      BinaryName,
		Usage:     "Compares two configuration files...",
		ArgsUsage: "<old_file> <new_file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format (stylish, plain, json)",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return fmt.Errorf("expected 2 file paths, got %d", cmd.NArg())
			}

			format := cmd.String("format")
			oldPath := cmd.Args().Get(0)
			newPath := cmd.Args().Get(1)

			diff, err := code.GenDiff(oldPath, newPath, format)
			if err != nil {
				return err
			}
			fmt.Print(diff)

			return nil
		},
	}

	return cmd.Run(context.Background(), args)
}

func Execute(args []string, stderr io.Writer) int {
	if err := Run(args); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitError
	}

	return ExitOK
}
