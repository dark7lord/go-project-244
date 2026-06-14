// Package diff_builder provides functions to generate a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure.
package diff_builder

import (
	"maps"
	"slices"

	"github.com/samber/lo"

	"code/internal/diff"
)

// ValueKind represents the runtime type of a parsed JSON/YAML value.
type ValueKind string

const (
	Num     ValueKind = "num"
	String  ValueKind = "string"
	Bool    ValueKind = "bool"
	Null    ValueKind = "null"
	Map     ValueKind = "map"
	Array   ValueKind = "array"
	Unknown ValueKind = "unknown type"
)

func valueKind(value any) ValueKind {
	switch value.(type) {
	case float64, int:
		return Num
	case string:
		return String
	case bool:
		return Bool
	case nil:
		return Null
	case map[string]any:
		return Map
	case []any:
		return Array
	default:
		return Unknown
	}
}

func isEqual(a, b any) bool {
	typeA := valueKind(a)
	typeB := valueKind(b)

	if typeA == Unknown || typeB == Unknown {
		return false
	}

	if typeA == Array && typeB == typeA {
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

	if typeA == Map && typeB == typeA {
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

	if typeA == Num {
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
	typeA := valueKind(dataA)
	typeB := valueKind(dataB)

	normA := normalizeValue(dataA)
	normB := normalizeValue(dataB)

	if typeA != typeB || !isEqual(normA, normB) {
		typeDiff = diff.Changed
	}

	return BuildDiff(typeDiff, normA, normB)
}
