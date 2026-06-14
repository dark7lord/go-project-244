package diff_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"code/internal/diff"
)

const (
	testKey    = "key"
	testNested = "nested"
	testArr    = "arr"
)

func nestedDiff(children ...diff.Node) diff.Node {
	return diff.Node{
		TypeDiff: diff.Nested,
		Children: children,
	}
}

func TestGenMapDiff(t *testing.T) {
	tests := []struct {
		name string
		mapA map[string]any
		mapB map[string]any
		want diff.Node
	}{
		{
			name: "flat add",
			mapA: map[string]any{},
			mapB: map[string]any{testKey: 42.0},
			want: nestedDiff(
				diff.Node{Key: testKey, TypeDiff: diff.Added, NewValue: 42.0},
			),
		},
		{
			name: "flat remove",
			mapA: map[string]any{testKey: "value"},
			mapB: map[string]any{},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Removed, OldValue: "value"},
			),
		},
		{
			name: "flat changed",
			mapA: map[string]any{testKey: true},
			mapB: map[string]any{testKey: false},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Changed, OldValue: true, NewValue: false},
			),
		},
		{
			name: "flat unchanged",
			mapA: map[string]any{testKey: nil},
			mapB: map[string]any{testKey: nil},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Unchanged, OldValue: nil},
			),
		},
		{
			name: "nested object",
			mapA: map[string]any{testNested: map[string]any{"a": 1.0}},
			mapB: map[string]any{testNested: map[string]any{"a": 1.0, "b": 2.0}},
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
			name: "nested array",
			mapA: map[string]any{testArr: []any{1.0, 2.0}},
			mapB: map[string]any{testArr: []any{1.0, 3.0, 4.0}},
			want: nestedDiff(
				diff.Node{
					Key:      testArr,
					TypeDiff: diff.Changed,
					OldValue: []any{1.0, 2.0},
					NewValue: []any{1.0, 3.0, 4.0},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genMapDiff(tt.mapA, tt.mapB)
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
			name: "unknown type",
			a:    struct{}{},
			b:    struct{}{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isEqual(tt.a, tt.b))
		})
	}
}

func TestTypeVarAndNormalizeValue(t *testing.T) {
	assert.Equal(t, "num", typeVar(1))
	assert.Equal(t, "num", typeVar(1.0))
	assert.Equal(t, "string", typeVar("text"))
	assert.Equal(t, "bool", typeVar(false))
	assert.Equal(t, "null", typeVar(nil))
	assert.Equal(t, "map", typeVar(map[string]any{}))
	assert.Equal(t, testArr, typeVar([]any{}))
	assert.Equal(t, unknownType, typeVar(struct{}{}))

	assert.Equal(t, float64(1), normalizeValue(1))
	assert.Equal(t, float64(1), normalizeValue(1.0))
}

func TestBuildDiff(t *testing.T) {
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Added, OldValue: nil, NewValue: 1.0},
		BuildDiff(diff.Added, 1, 1.0),
	)
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Changed, OldValue: 1.0, NewValue: 2.0},
		BuildDiff(diff.Changed, 1, 2),
	)
}

func TestRecursiveGendiff(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want any
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
			assert.Equal(t, tt.want, RecursiveGendiff(tt.a, tt.b))
		})
	}
}
