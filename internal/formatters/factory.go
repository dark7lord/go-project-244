package formatters

import (
	"fmt"

	"code/internal/diff"
)

// Formatter renders a diff node.
type Formatter interface {
	Format(diff.Node) (string, error)
}

type stylishFormatter struct{}

func (stylishFormatter) Format(df diff.Node) (string, error) {
	return renderStylish(df, 0), nil
}

type plainFormatter struct{}

func (plainFormatter) Format(df diff.Node) (string, error) {
	return renderPlain(df, []string{}), nil
}

type jsonFormatter struct{}

func (jsonFormatter) Format(df diff.Node) (string, error) {
	return renderJSON(df)
}

// NewFormatter returns a formatter for the provided output format.
func NewFormatter(format string) (Formatter, error) {
	switch OutputFormat(format) {
	case Stylish:
		return stylishFormatter{}, nil
	case Plain:
		return plainFormatter{}, nil
	case JSON:
		return jsonFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
