package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/diff"
)

func isMap(diff any) bool {
	_, isMap := diff.(map[string]any)
	return isMap
}

func printStylishDiff(difference any, deep int) string {
	if !isMap(difference) {
		if difference == nil {
			return "null"
		}

		return fmt.Sprintf("%v", difference)
	}

	var builder strings.Builder

	writeValue := func(prefix, key string, value any, deep int) {
		pad := strings.Repeat("  ", deep)
		nl := ""

		if !isMap(value) {
			nl = "\n"
		}

		strValue := printStylishDiff(value, deep+1)

		if key == "" {
			fmt.Fprintf(&builder, "%s%s %s%s", pad, prefix, strValue, nl)
		}

		fmt.Fprintf(&builder, "%s%s %s: %s%s", pad, prefix, key, strValue, nl)
	}

	pad := strings.Repeat("  ", deep)

	switch v := difference.(type) {
	case map[string]any:
		keys := slices.Collect(maps.Keys(v))
		slices.Sort(keys)
		builder.WriteString("{\n")
		for _, key := range keys {
			d, isDiff := v[key].(diff.Diff)

			if !isDiff {
				writeValue(" ", key, v[key], deep+1)
				continue
			}

			switch d.TypeDiff {
			case diff.DiffTypeUnchanged:
				writeValue(" ", key, d.OldValue, deep+1)
			case diff.DiffTypeAdded:
				writeValue("+", key, d.NewValue, deep+1)
			case diff.DiffTypeRemoved:
				writeValue("-", key, d.OldValue, deep+1)
			case diff.DiffTypeChanged:
				writeValue("-", key, d.OldValue, deep+1)
				writeValue("+", key, d.NewValue, deep+1)
			}
		}

		builder.WriteString(pad)
		builder.WriteString("}")
		if deep != 0 {
			builder.WriteString("\n")
		}

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
