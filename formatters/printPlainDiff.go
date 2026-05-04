package formatters

import (
	"fmt"
	"maps"
	"slices"
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

func printPlainDiff(difference any, keys []string) string {
	switch v := difference.(type) {
	case diff.Diff:
		strPath := strings.Join(keys, ".")
		switch v.TypeDiff {
		case diff.Added:
			strValue := formatValue(v.NewValue)
			return fmt.Sprintf("Property '%s' was added with value: %s", strPath, strValue)
		case diff.Removed:
			return fmt.Sprintf("Property '%s' was removed", strPath)
		case diff.Changed:
			strOldValue := formatValue(v.OldValue)
			strNewValue := formatValue(v.NewValue)
			return fmt.Sprintf("Property '%s' was updated. From %s to %s", strPath, strOldValue, strNewValue)
		default:
			return printPlainDiff(v.OldValue, keys) // for nested diffs
		}
	case map[string]any:
		lines := []string{}
		mapKeys := slices.Sorted(maps.Keys(v))
		for _, mk := range mapKeys {
			newPath := append(slices.Clone(keys), mk)
			if line := printPlainDiff(v[mk], newPath); line != "" {
				lines = append(lines, line)
			}
		}

		return strings.Join(lines, "\n")
	// TODO: подумать в тестах как выводит массивы в plain
	case []any:
		lines := []string{}
		for i, item := range v {
			strIndex := fmt.Sprintf("[%d]", i)
			newPath := append(slices.Clone(keys), strIndex)
			line := printPlainDiff(item, newPath)
			if line != "" {
				lines = append(lines, line)
			}
		}

		return strings.Join(lines, "\n")
	default:
		return ""
	}
}
