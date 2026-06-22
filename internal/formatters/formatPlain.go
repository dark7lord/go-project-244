package formatters

import (
	"fmt"
	"slices"
	"strings"

	"code/internal/diff"
)

func formatValue(value any) string {
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

func formatPlain(df diff.Node, keys []string) string {
	var strKeys string

	if df.TypeDiff != diff.Nested {
		strKeys = strings.Join(keys, ".")
	}

	switch df.TypeDiff {
	case diff.Added:
		strValue := formatValue(diff.ToNative(df.NewValue))
		return fmt.Sprintf("Property '%s' was added with value: %s", strKeys, strValue)
	case diff.Removed:
		return fmt.Sprintf("Property '%s' was removed", strKeys)
	case diff.Changed:
		strOldValue := formatValue(diff.ToNative(df.OldValue))
		strNewValue := formatValue(diff.ToNative(df.NewValue))
		return fmt.Sprintf("Property '%s' was updated. From %s to %s", strKeys, strOldValue, strNewValue)
	case diff.Nested:
		var builder strings.Builder

		if len(df.Children) == 0 {
			return ""
		}

		for _, child := range df.Children {
			childKeys := append(slices.Clone(keys), child.Key)
			result := formatPlain(child, childKeys)

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
