// Package code provides functions for getting the difference between files
package code

import (
	"code/diff"
	"code/formatters"
	"code/parsers"
)

// Types of differences between two files.
const (
	unknownType = "unknown type"
	arr         = "arr"
)

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
		return arr
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

	if typeA == arr && typeB == arr {
		sliceA := a.([]any)
		sliceB := b.([]any)

		if len(sliceA) != len(typeB) {
			return false
		}

		for i := range sliceA {
			if !isEqual(sliceA[i], sliceB[i]) {
				return false
			}
		}

		return true
	}

	if typeA != typeB {
		return false
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

func genMapDiff(mapA, mapB map[string]any) map[string]any {
	diffs := map[string]any{}
	for key, valueA := range mapA {
		valueB, ok := mapB[key]

		if !ok {
			difference := BuildDiff(diff.Removed, valueA, nil)
			difference.Key = key
			diffs[key] = difference

			continue
		}

		diffs[key] = recursiveGendiff(valueA, valueB)
	}

	for key, valueB := range mapB {
		if _, ok := mapA[key]; !ok {
			diffs[key] = diff.Diff{
				Key:      key,
				TypeDiff: diff.Added,
				NewValue: valueB,
			}
		}
	}

	return diffs
}

func genSliceDiff(sliceA, sliceB []any) []any {
	result := []any{}
	lenA := len(sliceA)
	lenB := len(sliceB)

	for i, itemA := range sliceA {
		if i < lenB {
			itemB := sliceB[i]
			diff := recursiveGendiff(itemA, itemB)
			result = append(result, diff)
		}
	}

	if lenA > lenB {
		for _, itemA := range sliceA[lenB:] {
			d := BuildDiff(diff.Removed, itemA, nil)
			result = append(result, d)
		}
	}

	if lenB > lenA {
		for _, itemB := range sliceB[lenA:] {
			d := BuildDiff(diff.Added, nil, itemB)
			result = append(result, d)
		}
	}

	return result
}

// BuildDiff function builds a Diff struct based on the type of difference and the old and new values
func BuildDiff(typeDiff diff.Kind, oldValue, newValue any) diff.Diff {
	result := diff.Diff{
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

func recursiveGendiff(dataA, dataB any) any {
	mapA, isAMap := dataA.(map[string]any)
	mapB, isBmap := dataB.(map[string]any)

	if isAMap && isBmap {
		return genMapDiff(mapA, mapB)
	}

	sliceA, isASlice := dataA.([]any)
	sliceB, isBSlice := dataB.([]any)

	if isASlice && isBSlice {
		return genSliceDiff(sliceA, sliceB)
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

// GenDiff function returns the difference between two structures as a string
func GenDiff(pathA, pathB, format string) (string, error) {
	dataA, err := parsers.Parse(pathA)
	if err != nil {
		return "", err
	}

	dataB, err := parsers.Parse(pathB)
	if err != nil {
		return "", err
	}

	diff := recursiveGendiff(dataA, dataB)
	typedFormat := formatters.PrintFormat(format)
	result, err := formatters.PrintDiff(diff, typedFormat)
	if err != nil {
		return "", err
	}

	return result, nil
}
