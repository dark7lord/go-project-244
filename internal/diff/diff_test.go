package diff

import (
	"encoding/json"
	"testing"
)

func TestNodeMarshalJSON_Added(t *testing.T) {
	node := Node{Key: "k", TypeDiff: Added, NewValue: "v"}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "added" {
		t.Errorf("expected type 'added', got %v", result["type"])
	}
	if result["newValue"] != "v" {
		t.Errorf("expected newValue 'v', got %v", result["newValue"])
	}
	if _, ok := result["oldValue"]; ok {
		t.Error("unexpected oldValue for added node")
	}
}

func TestNodeMarshalJSON_Removed(t *testing.T) {
	node := Node{Key: "k", TypeDiff: Removed, OldValue: 42.0}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "removed" {
		t.Errorf("expected type 'removed', got %v", result["type"])
	}
	if result["oldValue"] != 42.0 {
		t.Errorf("expected oldValue 42, got %v", result["oldValue"])
	}
	if _, ok := result["newValue"]; ok {
		t.Error("unexpected newValue for removed node")
	}
}

func TestNodeMarshalJSON_Changed(t *testing.T) {
	node := Node{Key: "k", TypeDiff: Changed, OldValue: "old", NewValue: "new"}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "changed" {
		t.Errorf("expected type 'changed', got %v", result["type"])
	}
	if result["oldValue"] != "old" {
		t.Errorf("expected oldValue 'old', got %v", result["oldValue"])
	}
	if result["newValue"] != "new" {
		t.Errorf("expected newValue 'new', got %v", result["newValue"])
	}
}

func TestNodeMarshalJSON_Unchanged(t *testing.T) {
	node := Node{Key: "k", TypeDiff: Unchanged, OldValue: true}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "unchanged" {
		t.Errorf("expected type 'unchanged', got %v", result["type"])
	}
	if result["oldValue"] != true {
		t.Errorf("expected oldValue true, got %v", result["oldValue"])
	}
	if _, ok := result["newValue"]; ok {
		t.Error("unexpected newValue for unchanged node")
	}
}

func TestNodeMarshalJSON_Nested(t *testing.T) {
	node := Node{
		TypeDiff: Nested,
		Children: []Node{
			{Key: "a", TypeDiff: Added, NewValue: 1.0},
		},
	}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "nested" {
		t.Errorf("expected type 'nested', got %v", result["type"])
	}
	children, ok := result["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatal("expected 1 child")
	}
}

func TestNodeMarshalJSON_PreservesNull(t *testing.T) {
	node := Node{
		TypeDiff: Nested,
		Children: []Node{
			{Key: "nullVal", TypeDiff: Added, NewValue: nil},
			{Key: "zeroVal", TypeDiff: Changed, OldValue: 0.0, NewValue: ""},
		},
	}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	children := result["children"].([]any)
	first := children[0].(map[string]any)
	if first["newValue"] != nil {
		t.Error("expected null newValue")
	}
	second := children[1].(map[string]any)
	if second["oldValue"] != 0.0 {
		t.Errorf("expected oldValue 0, got %v", second["oldValue"])
	}
	if second["newValue"] != "" {
		t.Errorf("expected newValue '', got %v", second["newValue"])
	}
}

func TestNodeMarshalJSON_EmptyNode(t *testing.T) {
	node := Node{}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["type"] != "" {
		t.Errorf("expected empty type, got %v", result["type"])
	}
}
