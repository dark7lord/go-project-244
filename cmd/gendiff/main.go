// Package main implements a CLI utility to get the difference between two files like git diff
package main

import (
	"code"
	"code/parser"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
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

			filedata1, err1 := parser.Parse(filepath1)
			filedata2, err2 := parser.Parse(filepath2)

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

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
