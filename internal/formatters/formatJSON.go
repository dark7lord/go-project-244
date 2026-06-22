package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

func formatJSON(df diff.Node) (string, error) {
	if df.TypeDiff != diff.Nested || len(df.Children) == 0 {
		return "{}", nil
	}

	result := buildJSONMap(df.Children)
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	return string(jsonBytes), nil
}

func buildJSONMap(children []diff.Node) map[string]any {
	result := make(map[string]any, len(children))

	for _, child := range children {
		if child.TypeDiff == diff.Nested {
			result[child.Key] = map[string]any{
				"type":     "nested",
				"children": buildJSONMap(child.Children),
			}
		} else if leaf := buildLeaf(child); leaf != nil {
			result[child.Key] = leaf
		}
	}

	return result
}

func buildLeaf(n diff.Node) map[string]any {
	switch n.TypeDiff {
	case diff.Added:
		return map[string]any{
			"type":  "added",
			"value": n.NewValue,
		}

	case diff.Removed:
		return map[string]any{
			"type":  "removed",
			"value": n.OldValue,
		}

	case diff.Changed:
		return map[string]any{
			"type":     "changed",
			"oldValue": n.OldValue,
			"newValue": n.NewValue,
		}

	case diff.Unchanged:
		return map[string]any{
			"type":  "unchanged",
			"value": n.OldValue,
		}
	}

	return nil
}
