// Package diff_builder provides functions to generate a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure.
package diff_builder

import (
	"maps"
	"slices"

	"github.com/samber/lo"

	"code/internal/diff"
)

// Types of differences between two files.
const unknownType = "unknown type"

func typeVar(value any) string {
	switch value.(type) {
	case float64, int:
		return "num"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "null"
	case map[string]any:
		return "map"
	case []any:
		return "arr"
	default:
		return unknownType
	}
}

func isEqual(a, b any) bool {
	typeA := typeVar(a)
	typeB := typeVar(b)

	if typeA == unknownType || typeB == unknownType {
		return false
	}

	if typeA == "arr" && typeB == typeA {
		sliceA := a.([]any)
		sliceB := b.([]any)

		if len(sliceA) != len(sliceB) {
			return false
		}

		for i := range sliceA {
			if !isEqual(sliceA[i], sliceB[i]) {
				return false
			}
		}

		return true
	}

	if typeA == "map" && typeB == typeA {
		mapA := a.(map[string]any)
		mapB := b.(map[string]any)

		if len(mapA) != len(mapB) {
			return false
		}

		for key, valueA := range mapA {
			valueB, ok := mapB[key]
			if !ok {
				return false
			}

			if !isEqual(valueA, valueB) {
				return false
			}
		}

		return true
	}

	if typeA != typeB {
		return false
	}

	if typeA == "num" {
		return normalizeValue(a) == normalizeValue(b)
	}

	return a == b
}

// yaml parser return int instead of float64
func normalizeValue(value any) any {
	switch v := value.(type) {
	case int:
		return float64(v)
	default:
		return value
	}
}

func genMapDiff(mapA, mapB map[string]any) diff.Node {
	leftKeys := slices.Collect(maps.Keys(mapA))
	rightKeys := slices.Collect(maps.Keys(mapB))

	removed, added := lo.Difference(leftKeys, rightKeys)
	common := lo.Union(leftKeys, rightKeys)
	slices.Sort(common)

	nodes := make([]diff.Node, len(common))

	for i, key := range common {
		node := diff.Node{
			Key:      key,
			TypeDiff: diff.Unchanged,
		}

		switch {
		case slices.Contains(removed, key):
			node.TypeDiff = diff.Removed
			node.OldValue = mapA[key]
		case slices.Contains(added, key):
			node.TypeDiff = diff.Added
			node.NewValue = mapB[key]
		default:
			leftValue := mapA[key]
			rightValue := mapB[key]

			node = RecursiveGendiff(leftValue, rightValue)
			node.Key = key

			if !isEqual(leftValue, rightValue) {
				typeDiff := diff.Changed

				if len(node.Children) > 0 {
					typeDiff = diff.Nested
				}

				node.TypeDiff = typeDiff
			}
		}

		nodes[i] = node
	}

	resultNode := diff.Node{
		TypeDiff: diff.Nested,
		Children: nodes,
	}

	return resultNode
}

// BuildDiff function builds a Diff struct based on the type of difference and the old and new values
func BuildDiff(typeDiff diff.Kind, oldValue, newValue any) diff.Node {
	result := diff.Node{
		TypeDiff: typeDiff,
		OldValue: normalizeValue(oldValue),
	}

	if typeDiff == diff.Added {
		result.OldValue = nil
	}

	if typeDiff == diff.Added || typeDiff == diff.Changed {
		result.NewValue = normalizeValue(newValue)
	}

	return result
}

// RecursiveGendiff function recursively generates a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure
func RecursiveGendiff(dataA, dataB any) diff.Node {
	mapA, isAMap := dataA.(map[string]any)
	mapB, isBmap := dataB.(map[string]any)

	if isAMap && isBmap {
		return genMapDiff(mapA, mapB)
	}

	typeDiff := diff.Unchanged
	typeA := typeVar(dataA)
	typeB := typeVar(dataB)

	normA := normalizeValue(dataA)
	normB := normalizeValue(dataB)

	if typeA != typeB || !isEqual(normA, normB) {
		typeDiff = diff.Changed
	}

	return BuildDiff(typeDiff, normA, normB)
}
