package formatters

import (
	"fmt"
	"slices"
	"strings"

	"code/internal/diff"
)

type stylishFormatter struct{}

func (stylishFormatter) Format(df diff.Node) (string, error) {
	return renderStylish(df, 0), nil
}

// NewStylishFormatter returns a formatter that renders diff in stylish format.
func NewStylishFormatter() Formatter {
	return stylishFormatter{}
}

func renderStylishSlice(slice []any, deep int) string {
	if len(slice) == 0 {
		return "[]"
	}

	var builder strings.Builder

	builder.WriteString("[\n")
	innerPad := strings.Repeat("    ", deep+1)
	pad := strings.Repeat(" ", deep*4)

	for _, item := range slice {
		builder.WriteString(innerPad)
		builder.WriteString(renderStylishValue(item, deep+1))
		builder.WriteByte('\n')
	}

	builder.WriteString(pad)
	builder.WriteByte(']')

	return builder.String()
}

func renderStylishMap(mp map[string]any, deep int) string {
	if len(mp) == 0 {
		return "{}"
	}

	var builder strings.Builder
	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	builder.WriteString("{\n")
	innerPad := strings.Repeat("    ", deep+1)
	pad := strings.Repeat(" ", deep*4)

	for _, key := range keys {
		childVal := mp[key]
		builder.WriteString(innerPad)
		fmt.Fprintf(&builder, "%s: ", key)
		builder.WriteString(renderStylishValue(childVal, deep+1))
		builder.WriteByte('\n')
	}

	builder.WriteString(pad)
	builder.WriteByte('}')

	return builder.String()
}

func renderStylishValue(value any, deep int) string {
	switch v := value.(type) {
	case []any:
		return renderStylishSlice(v, deep)
	case map[string]any:
		return renderStylishMap(v, deep)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func renderStylish(df diff.Node, deep int) string {
	var builder strings.Builder

	repeats := (4 * deep) - 2
	pad := strings.Repeat(" ", max(repeats, 0))

	writeKeyValue := func(prefix, key string, value any) {
		if deep == 0 {
			trimmedPrefix := strings.TrimLeft(prefix, " ")
			builder.WriteString(trimmedPrefix)
			builder.WriteString(renderStylishValue(value, deep))

			return
		}

		builder.WriteString(pad)
		builder.WriteString(prefix)
		fmt.Fprintf(&builder, "%s: %s", key, renderStylishValue(value, deep))
	}

	switch df.TypeDiff {
	case diff.Unchanged:
		writeKeyValue("  ", df.Key, df.OldValue)

	case diff.Removed:
		writeKeyValue("- ", df.Key, df.OldValue)

	case diff.Added:
		writeKeyValue("+ ", df.Key, df.NewValue)

	case diff.Changed:
		writeKeyValue("- ", df.Key, df.OldValue)
		builder.WriteByte('\n')
		writeKeyValue("+ ", df.Key, df.NewValue)

	case diff.Nested:
		if deep > 0 {
			builder.WriteString(pad)
			fmt.Fprintf(&builder, "  %s: ", df.Key)
		}

		if len(df.Children) == 0 {
			builder.WriteString("{}")
			return builder.String()
		}

		builder.WriteString("{")

		for _, child := range df.Children {
			builder.WriteByte('\n')
			childStr := renderStylish(child, deep+1)
			builder.WriteString(childStr)
		}

		builder.WriteByte('\n')
		builder.WriteString(strings.Repeat(" ", deep*4))
		builder.WriteString("}")
	}

	return builder.String()
}
