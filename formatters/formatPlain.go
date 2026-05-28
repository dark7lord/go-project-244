package formatters

import (
	"fmt"
	"slices"
	// "maps"
	// "slices"
	"strings"

	"code/diff"
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
		strValue := formatValue(df.NewValue)
		return fmt.Sprintf("Property '%s' was added with value: %s", strKeys, strValue)
	case diff.Removed:
		return fmt.Sprintf("Property '%s' was removed", strKeys)
	case diff.Changed:
		strOldValue := formatValue(df.OldValue)
		strNewValue := formatValue(df.NewValue)
		return fmt.Sprintf("Property '%s' was updated. From %s to %s", strKeys, strOldValue, strNewValue)
	case diff.Nested:
		var builder strings.Builder

		// if leftValue {} equals rigntValue {}
		if len(df.Children) == 0 {
			return ""
		}

		for _, child := range df.Children {
			childKeys := append(slices.Clone(keys), child.Key)
			result := formatPlain(child, childKeys)

			if len(result) == 0 {
				continue
			}

			builder.WriteString(result)

			if child.TypeDiff != diff.Nested {
				builder.WriteByte('\n')
			}
		}

		return builder.String()

	}

	return ""
}
