package formatters

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"code/internal/diff"
)

type addedEntry struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type removedEntry struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type changedEntry struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	OldValue any    `json:"oldValue"`
	NewValue any    `json:"newValue"`
}

type unchangedEntry struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

func flattenDiff(df diff.Node, path []string) []any {
	key := strings.Join(path, ".")

	if df.TypeDiff == diff.Nested {
		var entries []any
		for _, child := range df.Children {
			childPath := append(slices.Clone(path), child.Key)
			entries = append(entries, flattenDiff(child, childPath)...)
		}

		return entries
	}

	switch df.TypeDiff {
	case diff.Added:
		return []any{addedEntry{
			Key:   key,
			Type:  "added",
			Value: diff.ToNative(df.NewValue),
		}}

	case diff.Removed:
		return []any{removedEntry{
			Key:   key,
			Type:  "removed",
			Value: diff.ToNative(df.OldValue),
		}}

	case diff.Changed:
		return []any{changedEntry{
			Key:      key,
			Type:     "changed",
			OldValue: diff.ToNative(df.OldValue),
			NewValue: diff.ToNative(df.NewValue),
		}}

	case diff.Unchanged:
		return []any{unchangedEntry{
			Key:   key,
			Type:  "unchanged",
			Value: diff.ToNative(df.OldValue),
		}}
	}

	return nil
}

func formatJSON(df diff.Node) (string, error) {
	entries := flattenDiff(df, []string{})
	if len(entries) == 0 {
		return "[]", nil
	}

	jsonDiff, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	return string(jsonDiff), nil
}
