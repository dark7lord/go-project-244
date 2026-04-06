// Package code provides functions for getting the difference between files
package code

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

const (
	unknownType = "unknown type"
)

func typeVar(value any) string {
	switch value.(type) {
	case float64:
		return "num"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "null"
	case map[string]any:
		return "map"
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

	if typeA != typeB {
		return false
	}

	return a == b
}

// Diff struct represents the difference between two values
type Diff struct {
	typeDiff string
	oldValue any
	newValue any
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
		normA := normalizeValue(valueA)
		diff := Diff{
			typeDiff: "unchanged",
			oldValue: normA,
		}

		valueB, ok := mapB[key]

		if !ok {
			diff.typeDiff = "removed"
			diffs[key] = diff

			continue
		}

		normB := normalizeValue(valueB)

		if (isPrimitive(normA) || isPrimitive(normB)) && !isEqual(normA, normB) {
			diff.typeDiff = "changed"
			diff.newValue = valueB
			diffs[key] = diff

			continue
		}

		leftValue, isLeftMap := valueA.(map[string]any)
		rightValue, isRightMap := valueB.(map[string]any)

		if isLeftMap && isRightMap {
			diff.oldValue = genMapDiff(leftValue, rightValue)
		}

		diffs[key] = diff
	}

	for key, valueB := range mapB {
		if _, ok := mapA[key]; !ok {
			diffs[key] = Diff{
				typeDiff: "added",
				newValue: valueB,
			}
		}
	}

	return diffs
}

// GenDiff function returns the difference between two structures as a string
func GenDiff(dataA, dataB any) string {
	mapA := dataA.(map[string]any)
	mapB := dataB.(map[string]any)
	diff := genMapDiff(mapA, mapB)

	return printDiff(diff, 0)
}

func isPrimitive(value any) bool {
	switch value.(type) {
	case string, int, int64, float64, bool, nil:
		return true
	default:
		return false
	}
}

func printDiff(diff any, deep int) string {
	if isPrimitive(diff) {
		if diff == nil {
			return "null"
		}

		return fmt.Sprintf("%v", diff)
	}

	var builder strings.Builder

	writeValue := func(prefix, key string, value any, deep int) {
		pad := strings.Repeat("  ", deep)
		nl := ""
		if isPrimitive(value) {
			nl = "\n"
		}
		strValue := printDiff(value, deep+1)

		if key == "" {
			fmt.Fprintf(&builder, "%s%s %s%s", pad, prefix, strValue, nl)
		}

		fmt.Fprintf(&builder, "%s%s %s: %s%s", pad, prefix, key, strValue, nl)
	}

	pad := strings.Repeat("  ", deep)

	switch v := diff.(type) {
	case map[string]any:
		keys := slices.Collect(maps.Keys(v))
		slices.Sort(keys)
		builder.WriteString("{\n")
		for _, key := range keys {
			d, isDiff := v[key].(Diff)

			if !isDiff {
				writeValue(" ", key, v[key], deep+1)
				continue
			}

			switch d.typeDiff {
			case "unchanged":
				writeValue(" ", key, d.oldValue, deep+1)
			case "added":
				writeValue("+", key, d.newValue, deep+1)
			case "removed":
				writeValue("-", key, d.oldValue, deep+1)
			case "changed":
				writeValue("-", key, d.oldValue, deep+1)
				writeValue("+", key, d.newValue, deep+1)
			}
		}

		builder.WriteString(pad)
		builder.WriteString("}\n")
	case []any:
		if len(v) == 0 {
			builder.WriteString("[]\n")
		} else {
			builder.WriteString("[\n")
			for _, item := range v {
				writeValue(" ", "", item, deep+1)
			}
			builder.WriteString(pad)
			builder.WriteString("]\n")

		}
	}

	return builder.String()
}
