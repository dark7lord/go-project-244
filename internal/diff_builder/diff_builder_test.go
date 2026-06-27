package diff_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"code/internal/diff"
)

const (
	testKey      = "key"
	testNested   = "nested"
	testArrayKey = "arr"
)

func nestedDiff(children ...diff.Node) diff.Node {
	return diff.Node{
		TypeDiff: diff.Nested,
		Children: children,
	}
}

func TestBuildObjectDiff(t *testing.T) {
	tests := []struct {
		name  string
		left  map[string]any
		right map[string]any
		want  diff.Node
	}{
		{
			name:  "flat add",
			left:  map[string]any{},
			right: map[string]any{testKey: 42.0},
			want: nestedDiff(
				diff.Node{Key: testKey, TypeDiff: diff.Added, NewValue: 42.0},
			),
		},
		{
			name:  "flat remove",
			left:  map[string]any{testKey: "value"},
			right: map[string]any{},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Removed, OldValue: "value"},
			),
		},
		{
			name:  "flat changed",
			left:  map[string]any{testKey: true},
			right: map[string]any{testKey: false},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Changed, OldValue: true, NewValue: false},
			),
		},
		{
			name:  "flat unchanged",
			left:  map[string]any{testKey: nil},
			right: map[string]any{testKey: nil},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Unchanged, OldValue: nil},
			),
		},
		{
			name:  "nested object",
			left:  map[string]any{testNested: map[string]any{"a": 1.0}},
			right: map[string]any{testNested: map[string]any{"a": 1.0, "b": 2.0}},
			want: nestedDiff(
				diff.Node{
					Key:      testNested,
					TypeDiff: diff.Nested,
					Children: []diff.Node{
						{Key: "a", TypeDiff: diff.Unchanged, OldValue: 1.0},
						{Key: "b", TypeDiff: diff.Added, NewValue: 2.0},
					},
				},
			),
		},
		{
			name:  "nested array",
			left:  map[string]any{testArrayKey: []any{1.0, 2.0}},
			right: map[string]any{testArrayKey: []any{1.0, 3.0, 4.0}},
			want: nestedDiff(
				diff.Node{
					Key:      testArrayKey,
					TypeDiff: diff.Changed,
					OldValue: []any{1.0, 2.0},
					NewValue: []any{1.0, 3.0, 4.0},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildObjectDiff(tt.left, tt.right)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want bool
	}{
		{
			name: "same bool",
			a:    true,
			b:    true,
			want: true,
		},
		{
			name: "different bool",
			a:    true,
			b:    false,
			want: false,
		},
		{
			name: "different type",
			a:    1.0,
			b:    "1", want: false},
		{
			name: "same map",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"a": 1.0},
			want: true,
		},
		{
			name: "different map",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"a": 2.0},
			want: false,
		},
		{
			name: "maps same length different keys",
			a:    map[string]any{"a": 1.0, "b": 2.0},
			b:    map[string]any{"a": 1.0, "c": 3.0},
			want: false,
		},
		{
			name: "same slice",
			a:    []any{1.0, 2.0},
			b:    []any{1.0, 2.0},
			want: true,
		},
		{
			name: "different slice",
			a:    []any{1.0, 2.0},
			b:    []any{1.0, 3.0},
			want: false,
		},
		{
			name: "nil values",
			a:    nil,
			b:    nil,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isEqual(tt.a, tt.b))
		})
	}
}

func TestIsEqualUnreachable(t *testing.T) {
	assert.False(t, isEqual(struct{}{}, struct{}{}))
}

func TestBuildDiff(t *testing.T) {
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Added, OldValue: nil, NewValue: 1.0},
		buildDiff(diff.Added, 1.0, 1.0),
	)
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Changed, OldValue: 1.0, NewValue: 2.0},
		buildDiff(diff.Changed, 1.0, 2.0),
	)
}

func TestBuildDiffTree(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want diff.Node
	}{
		{
			name: "changed primitive",
			a:    true,
			b:    false,
			want: diff.Node{TypeDiff: diff.Changed, OldValue: true, NewValue: false},
		},
		{
			name: "unchanged nested map",
			a:    map[string]any{"k": "v"},
			b:    map[string]any{"k": "v"},
			want: nestedDiff(
				diff.Node{Key: "k", TypeDiff: diff.Unchanged, OldValue: "v"},
			),
		},
		{
			name: "changed nested slice",
			a:    []any{1.0, 2.0},
			b:    []any{1.0, 3.0},
			want: diff.Node{TypeDiff: diff.Changed, OldValue: []any{1.0, 2.0}, NewValue: []any{1.0, 3.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, BuildDiffTree(tt.a, tt.b))
		})
	}
}
