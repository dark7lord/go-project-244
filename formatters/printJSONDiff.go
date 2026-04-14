package formatters

import (
	"encoding/json"

	"code/diff"
)

func cleanDiff(difference any) any {
	switch d := difference.(type) {
	case diff.Diff:
		result := map[string]any{
			"TypeDiff": d.TypeDiff,
		}

		switch d.TypeDiff {
		case diff.DiffTypeAdded:
			result["NewValue"] = d.NewValue
		case diff.DiffTypeRemoved:
			result["OldValue"] = d.OldValue
		case diff.DiffTypeChanged:
			result["OldValue"] = d.OldValue
			result["NewValue"] = d.NewValue
		case diff.DiffTypeUnchanged:
			return cleanDiff(d.OldValue)
		}

		return result

	case map[string]any:
		result := make(map[string]any, len(d))
		for key, value := range d {
			result[key] = cleanDiff(value)
		}

		return result

	default:
		return difference
	}
}

func printJSONDiff(difference any) string {
	cleaned := cleanDiff(difference)
	jsonDiff, _ := json.MarshalIndent(cleaned, "", "  ")
	return string(jsonDiff)
}
