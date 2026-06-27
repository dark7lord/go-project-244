package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

func renderJSON(df diff.Node) (string, error) {
	if df.TypeDiff != diff.Nested || len(df.Children) == 0 {
		return "{}", nil
	}

	jsonBytes, err := json.MarshalIndent(df, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	return string(jsonBytes), nil
}
