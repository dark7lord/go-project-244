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
	"code/parsers"
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
				Usage:       "output format",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return errors.New("wrong number of args file paths")
			}

			filepath1 := cmd.Args().Get(0)
			filepath2 := cmd.Args().Get(1)

			filedata1, err1 := parsers.Parse(filepath1)
			filedata2, err2 := parsers.Parse(filepath2)

			if err1 != nil {
				return err1
			}

			if err2 != nil {
				return err2
			}

			result := code.GenDiff(filedata1, filedata2)
			fmt.Println(result)

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
