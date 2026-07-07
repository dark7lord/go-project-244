// Package diff_builder builds a difference tree between two parsed data structures.
package diff_builder

import (
	"reflect"
	"slices"

	"code/internal/diff"
)

func collectObjectKeys(left, right map[string]any) []string {
	allKeys := make([]string, 0, len(left)+len(right))

	for k := range left {
		allKeys = append(allKeys, k)
	}

	for k := range right {
		allKeys = append(allKeys, k)
	}

	slices.Sort(allKeys)

	return slices.Compact(allKeys)
}

func buildObjectDiff(left, right map[string]any) diff.Node {
	allKeys := collectObjectKeys(left, right)

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

			if !reflect.DeepEqual(leftValue, rightValue) {
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

	if !reflect.DeepEqual(left, right) {
		typeDiff = diff.Changed
	}

	return buildDiff(typeDiff, left, right)
}
