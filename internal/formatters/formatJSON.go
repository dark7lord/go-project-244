package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

func transformDiff(d diff.Node) any {
	result := map[string]any{}

	switch d.TypeDiff {
	case diff.Added:
		return diff.ToNative(d.NewValue)
	case diff.Removed:
		return diff.ToNative(d.OldValue)
	case diff.Changed:
		result["[old value]"] = diff.ToNative(d.OldValue)
		result["[new value]"] = diff.ToNative(d.NewValue)
	case diff.Unchanged:
		return diff.ToNative(d.OldValue)
	case diff.Nested:
		if len(d.Children) == 0 {
			return result
		}

		for _, child := range d.Children {
			var prefix string

			switch child.TypeDiff {
			case diff.Added:
				prefix = "added"
			case diff.Removed:
				prefix = "deleted"
			case diff.Changed:
				prefix = "changed"
			default:
				prefix = ""
			}

			resultKey := child.Key

			if prefix != "" {
				resultKey = fmt.Sprintf("%s [%s]", child.Key, prefix)
			}

			result[resultKey] = transformDiff(child)
		}
	}

	return result
}

func formatJSON(df diff.Node) (string, error) {
	cleaned := transformDiff(df)
	jsonDiff, err := json.MarshalIndent(cleaned, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	return string(jsonDiff), nil
}
