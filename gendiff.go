package code

import (
	"fmt"
	"slices"
	"strings"
)

func typeVar(variable any) string {
	switch variable.(type) {
	case float64:
		return "num"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "null"
	default:
		return "unknown type"
	}
}

func isEqual(a, b any) bool {
	if typeVar(a) != typeVar(b) {
		return false
	}

	return a == b
}

func GenDiff(dataA, dataB any) string {
	mapA := dataA.(map[string]any)
	mapB := dataB.(map[string]any)

	result := map[string]string{}
	keys := []string{}

	for key := range mapA {
		keys = append(keys, key)

		if _, ok := mapB[key]; !ok {
			result[key] = "removed"
			continue
		}

		if !isEqual(mapA[key], mapB[key]) {
			result[key] = "changed"
			continue
		}

		result[key] = "unchanged"
	}

	for key := range mapB {
		if _, ok := mapA[key]; !ok {
			result[key] = "added"
			keys = append(keys, key)
		}
	}

	pad := "  "
	var builder strings.Builder

	writeRow := func(prefix, key string, value any) {
		normalizedValue := value
		if value == nil {
			normalizedValue = "null"
		}
		builder.WriteString(pad)
		builder.WriteString(fmt.Sprintf("%s %s: %v", prefix, key, normalizedValue))
		builder.WriteRune('\n')
	}

	slices.Sort(keys)
	builder.WriteString("{\n")
	for _, key := range keys {
		switch result[key] {
		case "unchanged":
			writeRow(" ", key, mapA[key])
		case "added":
			writeRow("+", key, mapB[key])
		case "removed":
			writeRow("-", key, mapA[key])
		case "changed":
			writeRow("-", key, mapA[key])
			writeRow("+", key, mapB[key])
		}
	}
	builder.WriteString("}")
	return builder.String()
}
