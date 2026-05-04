package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/diff"
)

func isComplex(diff any) bool {
	_, isMap := diff.(map[string]any)
	_, isSlice := diff.([]any)
	return isMap || isSlice
}

func printStylishDiff(difference any, deep int) string {
	if !isComplex(difference) {
		if difference == nil {
			return "null"
		}

		return fmt.Sprintf("%v", difference)
	}

	var builder strings.Builder

	writeValue := func(prefix, key string, value any, deep int) {
		pad := strings.Repeat("  ", deep)
		nl := ""

		if !isComplex(value) {
			nl = "\n"
		}

		strValue := printStylishDiff(value, deep+1)

		if key == "" {
			fmt.Fprintf(&builder, "%s%s %s%s", pad, prefix, strValue, nl)
			return
		}

		fmt.Fprintf(&builder, "%s%s %s: %s%s", pad, prefix, key, strValue, nl)
	}

	printDiff := func(d diff.Diff, key string, deep int) {
		switch d.TypeDiff {
		case diff.Unchanged:
			writeValue(" ", key, d.OldValue, deep)
		case diff.Added:
			writeValue("+", key, d.NewValue, deep)
		case diff.Removed:
			writeValue("-", key, d.OldValue, deep)
		case diff.Changed:
			writeValue("-", key, d.OldValue, deep)
			writeValue("+", key, d.NewValue, deep)
		}
	}

	pad := strings.Repeat("  ", deep)

	switch v := difference.(type) {
	case map[string]any:
		keys := slices.Collect(maps.Keys(v))
		slices.Sort(keys)
		builder.WriteString("{\n")
		for _, key := range keys {
			d, isDiff := v[key].(diff.Diff)

			if isDiff {
				printDiff(d, key, deep+1)
				continue
			}

			writeValue(" ", key, v[key], deep+1)
		}

		builder.WriteString(pad)
		builder.WriteString("}")
		if deep != 0 {
			builder.WriteString("\n")
		}

	case []any:
		if len(v) == 0 {
			builder.WriteString("[]\n")
			return builder.String()
		}

		builder.WriteString("[\n")

		for _, item := range v {
			d, isDiff := item.(diff.Diff)
			if isDiff {
				printDiff(d, "", deep+1)
				continue
			}
			writeValue(" ", "", item, deep+1)
		}

		builder.WriteString(pad)
		builder.WriteString("]")
		if deep != 0 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
