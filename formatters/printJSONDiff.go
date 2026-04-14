package formatters

import (
	"encoding/json"

	"code"
)

func cleanDiff(diff any) any {
	switch d := diff.(type) {
	case code.Diff:
		result := map[string]any{
			"TypeDiff": d.TypeDiff,
		}

		switch d.TypeDiff {
		case code.DiffTypeAdded:
			result["NewValue"] = d.NewValue
		case code.DiffTypeRemoved:
			result["OldValue"] = d.OldValue
		case code.DiffTypeChanged:
			result["OldValue"] = d.OldValue
			result["NewValue"] = d.NewValue
		case code.DiffTypeUnchanged:
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
		return diff
	}
}

func printJSONDiff(diff any) string {
	cleaned := cleanDiff(diff)
	jsonDiff, _ := json.MarshalIndent(cleaned, "", "  ")
	return string(jsonDiff)
}
