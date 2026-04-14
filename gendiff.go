// Package code provides functions for getting the difference between files
package code

// Types of differences between two files.
const (
	unknownType = "unknown type"

	DiffTypeUnchanged = "unchanged"
	DiffTypeAdded     = "added"
	DiffTypeRemoved   = "removed"
	DiffTypeChanged   = "changed"
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
	TypeDiff string
	OldValue any
	NewValue any
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

// IsPrimitive function checks if the value is a primitive type (string, number, boolean, or null)
func IsPrimitive(value any) bool {
	switch value.(type) {
	case string, int, int64, float64, bool, nil:
		return true
	default:
		return false
	}
}

func genMapDiff(mapA, mapB map[string]any) map[string]any {
	diffs := map[string]any{}
	for key, valueA := range mapA {
		normA := normalizeValue(valueA)
		diff := Diff{
			TypeDiff: DiffTypeUnchanged,
			OldValue: valueA,
		}

		valueB, ok := mapB[key]

		if !ok {
			diff.TypeDiff = DiffTypeRemoved
			diffs[key] = diff

			continue
		}

		normB := normalizeValue(valueB)

		if (IsPrimitive(normA) || IsPrimitive(normB)) && !isEqual(normA, normB) {
			diff.TypeDiff = DiffTypeChanged
			diff.NewValue = valueB
			diffs[key] = diff

			continue
		}

		leftValue, isLeftMap := valueA.(map[string]any)
		rightValue, isRightMap := valueB.(map[string]any)

		if isLeftMap && isRightMap {
			diff.OldValue = genMapDiff(leftValue, rightValue)
		}

		diffs[key] = diff
	}

	for key, valueB := range mapB {
		if _, ok := mapA[key]; !ok {
			diffs[key] = Diff{
				TypeDiff: DiffTypeAdded,
				NewValue: valueB,
			}
		}
	}

	return diffs
}

// GenDiff function returns the difference between two structures as a string
func GenDiff(dataA, dataB any) map[string]any {
	mapA := dataA.(map[string]any)
	mapB := dataB.(map[string]any)
	diff := genMapDiff(mapA, mapB)

	return diff
}
