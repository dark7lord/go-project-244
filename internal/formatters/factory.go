package formatters

import (
	"fmt"

	"code/internal/diff"
)

// Formatter renders a diff node.
type Formatter interface {
	Format(diff.Node) (string, error)
}

// NewFormatter returns a formatter for the provided output format.
func NewFormatter(format string) (Formatter, error) {
	switch OutputFormat(format) {
	case Stylish:
		return NewStylishFormatter(), nil
	case Plain:
		return NewPlainFormatter(), nil
	case JSON:
		return NewJSONFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
