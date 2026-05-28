// Package code provides functions for getting the difference between files
package code

import (
	"code/diff_builder"
	"code/formatters"
	"code/parsers"
)

// GenDiff function returns the difference between two structures as a string
func GenDiff(pathA, pathB, format string) (string, error) {
	dataA, err := parsers.Parse(pathA)
	if err != nil {
		return "", err
	}

	dataB, err := parsers.Parse(pathB)
	if err != nil {
		return "", err
	}

	diff := diff_builder.RecursiveGendiff(dataA, dataB)
	typedFormat := formatters.PrintFormat(format)
	result, err := formatters.PrintDiff(diff, typedFormat)
	if err != nil {
		return "", err
	}

	return result, nil
}
