// Package formatters provides functions for formatting the difference between files
package formatters

import (
	"fmt"
	"slices"
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
func PrintDiff(diff any, format PrintFormat) (string, error) {
	switch format {
	case Stylish:
		return printStylishDiff(diff, 0), nil
	case JSON:
		return printJSONDiff(diff), nil
	case Plain:
		return printPlainDiff(diff, []string{}), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
