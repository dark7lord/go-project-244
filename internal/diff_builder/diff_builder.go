// Package diff_builder provides functions to generate a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure.
package diff_builder

import (
	"slices"

	"code/internal/diff"
)

func isEqual(a, b diff.Value) bool {
	switch va := a.(type) {
	case diff.Slice:
		vb, ok := b.(diff.Slice)
		if !ok || len(va) != len(vb) {
			return false
		}
		for i := range va {
			if !isEqual(va[i], vb[i]) {
				return false
			}
		}

		return true

	case diff.Map:
		vb, ok := b.(diff.Map)
		if !ok || len(va) != len(vb) {
			return false
		}
		for key, valA := range va {
			valB, exists := vb[key]
			if !exists || !isEqual(valA, valB) {
				return false
			}
		}

		return true

	case diff.Number:
		vb, ok := b.(diff.Number)
		return ok && va == vb

	case diff.String:
		vb, ok := b.(diff.String)
		return ok && va == vb

	case diff.Boolean:
		vb, ok := b.(diff.Boolean)
		return ok && va == vb

	case diff.Null:
		_, ok := b.(diff.Null)
		return ok

	default:
		return false
	}
}

func getDiffKeys(oldMap, newMap diff.Map) (removed, added, common []string) {
	lenA, lenB := len(oldMap), len(newMap)
	removed = make([]string, 0, lenA)
	added = make([]string, 0, lenB)
	common = make([]string, 0, max(lenA, lenB))

	for k := range oldMap {
		if _, ok := newMap[k]; ok {
			common = append(common, k)
		} else {
			removed = append(removed, k)
		}
	}

	slices.Sort(common)

	for k := range newMap {
		if _, ok := oldMap[k]; !ok {
			added = append(added, k)
		}
	}

	return removed, added, common
}

func genMapDiff(mapA, mapB diff.Map) diff.Node {
	removed, added, common := getDiffKeys(mapA, mapB)
	allKeys := append(append(removed, added...), common...)
	slices.Sort(allKeys)

	nodes := make([]diff.Node, 0, len(allKeys))

	for _, key := range allKeys {
		_, inA := mapA[key]
		_, inB := mapB[key]

		switch {
		case inA && !inB:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Removed,
				OldValue: mapA[key],
			})
		case !inA && inB:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Added,
				NewValue: mapB[key],
			})
		default:
			leftValue := mapA[key]
			rightValue := mapB[key]

			node := RecursiveGendiff(leftValue, rightValue)
			node.Key = key

			if !isEqual(leftValue, rightValue) {
				typeDiff := diff.Changed

				if len(node.Children) > 0 {
					typeDiff = diff.Nested
				}

				node.TypeDiff = typeDiff
			}

			nodes = append(nodes, node)
		}
	}

	return diff.Node{
		TypeDiff: diff.Nested,
		Children: nodes,
	}
}

// BuildDiff function builds a Diff struct based on the type of difference and the old and new values
func BuildDiff(typeDiff diff.Kind, oldValue, newValue diff.Value) diff.Node {
	result := diff.Node{
		TypeDiff: typeDiff,
		OldValue: oldValue,
	}

	if typeDiff == diff.Added {
		result.OldValue = nil
	}

	if typeDiff == diff.Added || typeDiff == diff.Changed {
		result.NewValue = newValue
	}

	return result
}

// RecursiveGendiff function recursively generates a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure
func RecursiveGendiff(dataA, dataB diff.Value) diff.Node {
	mapA, isAMap := dataA.(diff.Map)
	mapB, isBmap := dataB.(diff.Map)

	if isAMap && isBmap {
		return genMapDiff(mapA, mapB)
	}

	typeDiff := diff.Unchanged

	if !isEqual(dataA, dataB) {
		typeDiff = diff.Changed
	}

	return BuildDiff(typeDiff, dataA, dataB)
}
