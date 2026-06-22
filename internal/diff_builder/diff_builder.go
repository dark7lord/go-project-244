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
	}

	return false
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

	for k := range newMap {
		if _, ok := oldMap[k]; !ok {
			added = append(added, k)
		}
	}

	return removed, added, common
}

func genMapDiff(oldMap, newMap diff.Map) diff.Node {
	removed, added, common := getDiffKeys(oldMap, newMap)
	allKeys := append(append(removed, added...), common...)
	slices.Sort(allKeys)

	nodes := make([]diff.Node, 0, len(allKeys))

	for _, key := range allKeys {
		_, inOld := oldMap[key]
		_, inNew := newMap[key]

		switch {
		case inOld && !inNew:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Removed,
				OldValue: diff.ToNative(oldMap[key]),
			})
		case !inOld && inNew:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Added,
				NewValue: diff.ToNative(newMap[key]),
			})
		default:
			oldVal := oldMap[key]
			newVal := newMap[key]

			node := RecursiveGendiff(oldVal, newVal)
			node.Key = key

			if !isEqual(oldVal, newVal) {
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

func buildDiff(typeDiff diff.Kind, oldValue, newValue diff.Value) diff.Node {
	result := diff.Node{
		TypeDiff: typeDiff,
		OldValue: diff.ToNative(oldValue),
	}

	if typeDiff == diff.Added {
		result.OldValue = nil
	}

	if typeDiff == diff.Added || typeDiff == diff.Changed {
		result.NewValue = diff.ToNative(newValue)
	}

	return result
}

// RecursiveGendiff function recursively generates a diff between two data structures,
// which can be maps or slices, and returns the resulting diff structure
func RecursiveGendiff(dataA, dataB diff.Value) diff.Node {
	oldVal, isOldMap := dataA.(diff.Map)
	newVal, isNewMap := dataB.(diff.Map)

	if isOldMap && isNewMap {
		return genMapDiff(oldVal, newVal)
	}

	typeDiff := diff.Unchanged

	if !isEqual(dataA, dataB) {
		typeDiff = diff.Changed
	}

	return buildDiff(typeDiff, dataA, dataB)
}
