// Package formatters provides functions for formatting the difference between files
package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code"
)

// PrintDiff function returns the difference between two structures as a string in the specified format
func PrintDiff(diff any, format string) string {
	if format == "stylish" {
		return printStylishDiff(diff, 0)
	}

	return printPlainDiff(diff, []string{})
}

func printStylishDiff(diff any, deep int) string {
	if code.IsPrimitive(diff) {
		if diff == nil {
			return "null"
		}

		return fmt.Sprintf("%v", diff)
	}

	var builder strings.Builder

	writeValue := func(prefix, key string, value any, deep int) {
		pad := strings.Repeat("  ", deep)
		nl := ""
		if code.IsPrimitive(value) {
			nl = "\n"
		}
		strValue := printStylishDiff(value, deep+1)

		if key == "" {
			fmt.Fprintf(&builder, "%s%s %s%s", pad, prefix, strValue, nl)
		}

		fmt.Fprintf(&builder, "%s%s %s: %s%s", pad, prefix, key, strValue, nl)
	}

	pad := strings.Repeat("  ", deep)

	switch v := diff.(type) {
	case map[string]any:
		keys := slices.Collect(maps.Keys(v))
		slices.Sort(keys)
		builder.WriteString("{\n")
		for _, key := range keys {
			d, isDiff := v[key].(code.Diff)

			if !isDiff {
				writeValue(" ", key, v[key], deep+1)
				continue
			}

			switch d.TypeDiff {
			case "unchanged":
				writeValue(" ", key, d.OldValue, deep+1)
			case "added":
				writeValue("+", key, d.NewValue, deep+1)
			case "removed":
				writeValue("-", key, d.OldValue, deep+1)
			case "changed":
				writeValue("-", key, d.OldValue, deep+1)
				writeValue("+", key, d.NewValue, deep+1)
			}
		}

		builder.WriteString(pad)
		builder.WriteString("}\n")
	case []any:
		if len(v) == 0 {
			builder.WriteString("[]\n")
		} else {
			builder.WriteString("[\n")
			for _, item := range v {
				writeValue(" ", "", item, deep+1)
			}
			builder.WriteString(pad)
			builder.WriteString("]\n")

		}
	}

	return builder.String()
}

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

func printPlainDiff(diff any, path []string) string {
	switch v := diff.(type) {
	case code.Diff:
		strPath := strings.Join(path, ".")
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
			return printPlainDiff(v.OldValue, path) // for nested diffs
		}
	case map[string]any:
		var builder strings.Builder

		keys := slices.Collect(maps.Keys(v))
		slices.Sort(keys)

		for i, key := range keys {
			newPath := append(slices.Clone(path), key)
			line := printPlainDiff(v[key], newPath)

			if line == "" {
				continue
			}

			builder.WriteString(line)

			if i != len(keys)-1 {
				builder.WriteString("\n")
			}
		}

		return builder.String()
	default:
		return ""
	}
}
