// Package main implements a CLI utility to get the difference between two files like git diff
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"code"
)

func Run() error {
	cmd := &cli.Command{
		Name:                  "gendiff",
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
			if cmd.NArg() != 2 {
				return errors.New("wrong number of args file paths")
			}

			format := cmd.String("format")
			if format != "stylish" && format != "plain" && format != "json" {
				return fmt.Errorf("unsupported format: %s", format)
			}

			filepathA := cmd.Args().Get(0)
			filepathB := cmd.Args().Get(1)

			diff, err := code.GenDiff(filepathA, filepathB, format)
			if err != nil {
				return err
			}

			fmt.Println(diff)

			return nil
		},
	}

	return cmd.Run(context.Background(), os.Args)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
