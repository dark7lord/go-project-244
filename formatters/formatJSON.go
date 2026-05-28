package formatters

import (
	"encoding/json"
	"fmt"

	"code/diff"
)

func transformDiff(d diff.Node) any {
	result := map[string]any{}

	switch d.TypeDiff {
	case diff.Added:
		return d.NewValue
	case diff.Removed:
		return d.OldValue
	case diff.Changed:
		result["[old value]"] = d.OldValue
		result["[new value]"] = d.NewValue
	case diff.Unchanged:
		return d.OldValue
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

	// fmt.Println(df)

	if err != nil {
		return "", err
	}

	return string(jsonDiff), nil
}
