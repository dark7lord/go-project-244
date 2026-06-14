// Package formatters provides functions for formatting the difference between files
package formatters

import (
	"fmt"
	"slices"

	"code/internal/diff"
)

// PrintFormat type represents the format in which the difference between files can be printed
type PrintFormat string

// Constants of formats in which the difference between files can be printed
const (
	Stylish PrintFormat = "stylish"
	Plain   PrintFormat = "plain"
	JSON    PrintFormat = "json"
)

var formats = []PrintFormat{Stylish, Plain, JSON}

// IsValidFormat checks if the provided format is valid
func IsValidFormat(format string) bool {
	return slices.Contains(formats, PrintFormat(format))
}

// PrintDiff function returns the difference between two structures as a string in the specified format
func PrintDiff(diff diff.Node, format PrintFormat) (string, error) {
	switch format {
	case Stylish:
		return formatStylish(diff, 0), nil
	case JSON:
		return formatJSON(diff)
	case Plain:
		return formatPlain(diff, []string{}), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
