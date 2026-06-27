// Package diff provided common type Diff for all packages
package diff

import "encoding/json"

// Kind type represents the type of difference between two values
type Kind string

// Constants of differences between two values
const (
	Added     Kind = "added"
	Removed   Kind = "removed"
	Changed   Kind = "changed"
	Unchanged Kind = "unchanged"
	Nested    Kind = "nested"
)

// Node struct represents the difference between two values
type Node struct {
	Key      string `json:"key,omitempty"`
	TypeDiff Kind   `json:"type"`
	OldValue any    `json:"oldValue,omitempty"`
	NewValue any    `json:"newValue,omitempty"`
	Children []Node `json:"children,omitempty"`
}

// MarshalJSON preserves explicit null values for fields that belong to the node kind.
func (n Node) MarshalJSON() ([]byte, error) {
	type jsonNode struct {
		Key      string `json:"key,omitempty"`
		TypeDiff Kind   `json:"type"`
		OldValue *any   `json:"oldValue,omitempty"`
		NewValue *any   `json:"newValue,omitempty"`
		Children []Node `json:"children,omitempty"`
	}

	node := jsonNode{
		Key:      n.Key,
		TypeDiff: n.TypeDiff,
		Children: n.Children,
	}

	switch n.TypeDiff {
	case Added:
		node.NewValue = valuePtr(n.NewValue)
	case Removed, Unchanged:
		node.OldValue = valuePtr(n.OldValue)
	case Changed:
		node.OldValue = valuePtr(n.OldValue)
		node.NewValue = valuePtr(n.NewValue)
	}

	return json.Marshal(node)
}

func valuePtr(v any) *any {
	return &v
}
