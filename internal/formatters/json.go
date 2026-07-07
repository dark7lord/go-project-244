package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

type jsonFormatter struct{}

func (jsonFormatter) Format(df diff.Node) (string, error) {
	return renderJSON(df)
}

// NewJSONFormatter returns a formatter that renders diff in JSON format.
func NewJSONFormatter() Formatter {
	return jsonFormatter{}
}

func renderJSON(df diff.Node) (string, error) {
	if df.TypeDiff != diff.Nested || len(df.Children) == 0 {
		return "{}", nil
	}

	jsonBytes, err := json.MarshalIndent(convertNode(df), "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	return string(jsonBytes), nil
}

type jsonNode struct {
	Key      string     `json:"key,omitempty"`
	TypeDiff diff.Kind  `json:"type"`
	OldValue *any       `json:"oldValue,omitempty"`
	NewValue *any       `json:"newValue,omitempty"`
	Children []jsonNode `json:"children,omitempty"`
}

func convertNode(n diff.Node) jsonNode {
	node := jsonNode{
		Key:      n.Key,
		TypeDiff: n.TypeDiff,
	}

	for _, child := range n.Children {
		node.Children = append(node.Children, convertNode(child))
	}

	switch n.TypeDiff {
	case diff.Added:
		node.NewValue = &n.NewValue
	case diff.Removed, diff.Unchanged:
		node.OldValue = &n.OldValue
	case diff.Changed:
		node.OldValue = &n.OldValue
		node.NewValue = &n.NewValue
	}

	return node
}
