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

func getDiffKeys(left, right map[string]any) (removed, added, common []string) {
	leftLen, rightLen := len(left), len(right)
	removed = make([]string, 0, leftLen)
	added = make([]string, 0, rightLen)
	common = make([]string, 0, max(leftLen, rightLen))

	for k := range left {
		if _, ok := right[k]; ok {
			common = append(common, k)
		} else {
			removed = append(removed, k)
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			added = append(added, k)
		}
	}

	return removed, added, common
}

func buildObjectDiff(left, right map[string]any) diff.Node {
	removed, added, common := getDiffKeys(left, right)
	allKeys := append(append(removed, added...), common...)
	slices.Sort(allKeys)

	nodes := make([]diff.Node, 0, len(allKeys))

	for _, key := range allKeys {
		_, isInLeft := left[key]
		_, isInRight := right[key]

		switch {
		case isInLeft && !isInRight:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Removed,
				OldValue: left[key],
			})
		case !isInLeft && isInRight:
			nodes = append(nodes, diff.Node{
				Key:      key,
				TypeDiff: diff.Added,
				NewValue: right[key],
			})
		default:
			leftValue := left[key]
			rightValue := right[key]

			node := BuildDiffTree(leftValue, rightValue)
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

// BuildDiffTree builds a diff tree between two parsed data structures.
func BuildDiffTree(left, right any) diff.Node {
	leftObject, isLeftObject := left.(map[string]any)
	rightObject, isRightObject := right.(map[string]any)

	if isLeftObject && isRightObject {
		return buildObjectDiff(leftObject, rightObject)
	}

	typeDiff := diff.Unchanged

	if !isEqual(left, right) {
		typeDiff = diff.Changed
	}

	return buildDiff(typeDiff, left, right)
}
