// Package diff_builder builds a difference tree between two parsed data structures.
package diff_builder

import (
	"slices"

	"code/internal/diff"
)

func isEqual(a, b any) bool {
	switch va := a.(type) {
	case []any:
		vb, ok := b.([]any)
		if !ok || len(va) != len(vb) {
			return false
		}
		for i := range va {
			if !isEqual(va[i], vb[i]) {
				return false
			}
		}

		return true

	case map[string]any:
		vb, ok := b.(map[string]any)
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

	case float64:
		vb, ok := b.(float64)
		return ok && va == vb

	case string:
		vb, ok := b.(string)
		return ok && va == vb

	case bool:
		vb, ok := b.(bool)
		return ok && va == vb

	case nil:
		return b == nil
	}

	return false
}

func getDiffKeys(oldMap, newMap map[string]any) (removed, added, common []string) {
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

func genMapDiff(oldMap, newMap map[string]any) diff.Node {
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
				OldValue: oldMap[key],
			})
		case !inOld && inNew:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Added,
				NewValue: newMap[key],
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

func buildDiff(typeDiff diff.Kind, oldValue, newValue any) diff.Node {
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
func RecursiveGendiff(oldData, newData any) diff.Node {
	oldVal, isOldMap := oldData.(map[string]any)
	newVal, isNewMap := newData.(map[string]any)

	if isOldMap && isNewMap {
		return genMapDiff(oldVal, newVal)
	}

	typeDiff := diff.Unchanged

	if !isEqual(oldData, newData) {
		typeDiff = diff.Changed
	}

	return buildDiff(typeDiff, oldData, newData)
}
