package formatters

import (
	"encoding/json"

	dff "code/diff"
)

func cleanDiff(diff any) any {
	switch d := diff.(type) {
	case dff.Diff:
		result := map[string]any{
			"TypeDiff": d.TypeDiff,
		}

		switch d.TypeDiff {
		case "added":
			result["NewValue"] = d.NewValue
		case "removed":
			result["OldValue"] = d.OldValue
		case "changed":
			result["OldValue"] = d.OldValue
			result["NewValue"] = d.NewValue
		case "unchanged":
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
