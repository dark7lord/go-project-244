package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	dff "code/diff"
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

func printPlainDiff(diff any, keys []string) string {
	switch v := diff.(type) {
	case dff.Diff:
		strPath := strings.Join(keys, ".")
		switch v.TypeDiff {
		case "added":
			strValue := formatValue(v.NewValue)
			return fmt.Sprintf("Property '%s' was added with value: %s", strPath, strValue)
		case "removed":
			return fmt.Sprintf("Property '%s' was removed", strPath)
		case "changed":
			strOldValue := formatValue(v.OldValue)
			strNewValue := formatValue(v.NewValue)
			return fmt.Sprintf("Property '%s' was updated. From %s to %s", strPath, strOldValue, strNewValue)
		default:
			return printPlainDiff(v.OldValue, keys) // for nested diffs
		}
	case map[string]any:
		var builder strings.Builder

		mapKeys := slices.Collect(maps.Keys(v))
		slices.Sort(mapKeys)

		for i, mk := range mapKeys {
			newPath := append(slices.Clone(keys), mk)
			line := printPlainDiff(v[mk], newPath)

			if line == "" {
				continue
			}

			builder.WriteString(line)

			if i != len(mapKeys)-1 {
				builder.WriteString("\n")
			}
		}

		return builder.String()
	default:
		return ""
	}
}
