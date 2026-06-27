// Package formatters provides functions for formatting the difference between files.
package formatters

import "code/internal/diff"

// OutputFormat represents a supported diff output format.
type OutputFormat string

// Supported output formats.
const (
	Stylish OutputFormat = "stylish"
	Plain   OutputFormat = "plain"
	JSON    OutputFormat = "json"
)

// IsValidFormat reports whether the provided format is supported.
func IsValidFormat(format string) bool {
	_, err := NewFormatter(format)

	return err == nil
}

// Format returns the diff rendered in the specified output format.
func Format(format string, df diff.Node) (string, error) {
	formatter, err := NewFormatter(format)
	if err != nil {
		return "", err
	}

	return formatter.Format(df)
}
