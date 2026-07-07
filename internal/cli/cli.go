// Package cli provides command-line execution for gendiff.
package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"code"
)

// BinaryName is the command name used by the CLI.
const BinaryName = "gendiff"

// Exit codes returned by Execute.
const (
	ExitOK    = 0
	ExitError = 1
)

// NewCommand builds and configures the CLI command.
func NewCommand() *cli.Command {
	return &cli.Command{
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
}
