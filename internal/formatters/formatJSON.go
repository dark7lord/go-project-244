package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

type jsonKey = string

const (
	keyType     jsonKey = "type"
	keyValue    jsonKey = "value"
	keyOldValue jsonKey = "oldValue"
	keyNewValue jsonKey = "newValue"
	keyChildren jsonKey = "children"
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
				keyType:     "nested",
				keyChildren: buildJSONMap(child.Children),
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
			keyType:  "added",
			keyValue: n.NewValue,
		}

	case diff.Removed:
		return map[string]any{
			keyType:  "removed",
			keyValue: n.OldValue,
		}

	case diff.Changed:
		return map[string]any{
			keyType:     "changed",
			keyOldValue: n.OldValue,
			keyNewValue: n.NewValue,
		}

	case diff.Unchanged:
		return map[string]any{
			keyType:  "unchanged",
			keyValue: n.OldValue,
		}
	}

	return nil
}
