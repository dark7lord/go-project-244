package formatters

import (
	"fmt"
	"slices"
	"strings"

	"code/internal/diff"
)

func strSlice(slice []any, deep int) string {
	if len(slice) == 0 {
		return "[]"
	}

	var builder strings.Builder

	builder.WriteString("[\n")
	innerPad := strings.Repeat("    ", deep+1)
	pad := strings.Repeat(" ", deep*4)

	for _, item := range slice {
		builder.WriteString(innerPad)
		builder.WriteString(strValue(item, deep+1))
		builder.WriteByte('\n')
	}

	builder.WriteString(pad)
	builder.WriteByte(']')

	return builder.String()
}

func strMap(mp map[string]any, deep int) string {
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
		builder.WriteString(strValue(childVal, deep+1))
		builder.WriteByte('\n')
	}

	builder.WriteString(pad)
	builder.WriteByte('}')

	return builder.String()
}

func strValue(value any, deep int) string {
	switch v := value.(type) {
	case []any:
		return strSlice(v, deep)
	case map[string]any:
		return strMap(v, deep)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatStylish(df diff.Node, deep int) string {
	var builder strings.Builder

	repeats := (4 * deep) - 2
	pad := strings.Repeat(" ", max(repeats, 0))

	writeKeyValue := func(prefix, key string, value any) {
		if deep == 0 {
			trimmedPrefix := strings.TrimLeft(prefix, " ")
			builder.WriteString(trimmedPrefix)
			builder.WriteString(strValue(value, deep))

			return
		}

		builder.WriteString(pad)
		builder.WriteString(prefix)
		fmt.Fprintf(&builder, "%s: %s", key, strValue(value, deep))
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
			childStr := formatStylish(child, deep+1)
			builder.WriteString(childStr)
		}

		builder.WriteByte('\n')
		builder.WriteString(strings.Repeat(" ", deep*4))
		builder.WriteString("}")
	}

	return builder.String()
}
