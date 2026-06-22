// Package code provides functions for comparing two configuration files and formatting the difference
package code

import (
	"fmt"

	"code/internal/diff_builder"
	"code/internal/formatters"
	"code/internal/parsers"
)

// GenDiff function returns the difference between two structures as a string
func GenDiff(pathA, pathB, format string) (string, error) {
	if !formatters.IsValidFormat(format) {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

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

	return formatters.PrintDiff(diff, typedFormat)
}
