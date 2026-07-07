package diff

import "testing"

func TestNodeZeroValue(t *testing.T) {
	node := Node{}
	if node.Key != "" || node.TypeDiff != "" || node.OldValue != nil || node.NewValue != nil || len(node.Children) != 0 {
		t.Fatalf("unexpected zero value node: %#v", node)
	}
}
