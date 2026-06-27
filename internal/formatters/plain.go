package formatters

import (
	"fmt"
	"slices"
	"strings"

	"code/internal/diff"
)

func renderPlainValue(value any) string {
	switch value.(type) {
	case string:
		return fmt.Sprintf("'%s'", value)
	case nil:
		return "null"
	case []any, map[string]any:
		return "[complex value]"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func renderPlain(df diff.Node, keys []string) string {
	var propertyPath string

	if df.TypeDiff != diff.Nested {
		propertyPath = strings.Join(keys, ".")
	}

	switch df.TypeDiff {
	case diff.Added:
		value := renderPlainValue(df.NewValue)
		return fmt.Sprintf("Property '%s' was added with value: %s", propertyPath, value)
	case diff.Removed:
		return fmt.Sprintf("Property '%s' was removed", propertyPath)
	case diff.Changed:
		oldValue := renderPlainValue(df.OldValue)
		newValue := renderPlainValue(df.NewValue)
		return fmt.Sprintf("Property '%s' was updated. From %s to %s", propertyPath, oldValue, newValue)
	case diff.Nested:
		var builder strings.Builder

		if len(df.Children) == 0 {
			return ""
		}

		for _, child := range df.Children {
			childKeys := append(slices.Clone(keys), child.Key)
			result := renderPlain(child, childKeys)

			if len(result) == 0 {
				continue
			}

			if builder.Len() > 0 {
				builder.WriteByte('\n')
			}

			builder.WriteString(result)
		}

		return builder.String()
	}

	return ""
}
