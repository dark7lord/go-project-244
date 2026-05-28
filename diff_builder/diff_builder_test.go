package diff_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"code/diff"
)

func TestGenMapDiff(t *testing.T) {
	tests := []struct {
		name string
		mapA map[string]any
		mapB map[string]any
		want diff.Node
	}{
		{
			name: "flat add num",
			mapA: map[string]any{},
			mapB: map[string]any{"key": 42.0},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key",
						TypeDiff: diff.Added,
						NewValue: 42.0,
					},
				},
			},
		},
		{
			name: "flat remove string",
			mapA: map[string]any{"key": "value"},
			mapB: map[string]any{},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key",
						TypeDiff: diff.Removed,
						OldValue: "value",
					},
				},
			},
		},
		{
			name: "flat changed bool",
			mapA: map[string]any{"key": true},
			mapB: map[string]any{"key": false},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key",
						TypeDiff: diff.Changed,
						OldValue: true,
						NewValue: false,
					},
				},
			},
		},
		{
			name: "flat unchanged null",
			mapA: map[string]any{"key": nil},
			mapB: map[string]any{"key": nil},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key",
						TypeDiff: diff.Unchanged,
						OldValue: nil,
					},
				},
			},
		},
		{
			name: "nested objects",
			mapA: map[string]any{
				"nested": map[string]any{"a": 1.0},
			},
			mapB: map[string]any{
				"nested": map[string]any{"a": 1.0, "b": 2.0},
			},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "nested",
						TypeDiff: diff.Nested,
						Children: []diff.Node{
							{
								Key:      "a",
								TypeDiff: diff.Unchanged,
								OldValue: 1.0,
							},
							{
								Key:      "b",
								TypeDiff: diff.Added,
								NewValue: 2.0,
							},
						},
					},
				},
			},
		},
		{
			name: "nested arrays",
			mapA: map[string]any{
				"arr": []any{1.0, 2.0},
			},
			mapB: map[string]any{
				"arr": []any{1.0, 3.0, 4.0},
			},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "arr",
						TypeDiff: diff.Changed,
						OldValue: []any{1.0, 2.0},
						NewValue: []any{1.0, 3.0, 4.0},
					},
				},
			},
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
			name: "same number",
			a:    1.0,
			b:    1.0,
			want: true,
		},
		{
			name: "different type",
			a:    1.0,
			b:    "1",
			want: false,
		},
		{
			name: "same map",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"a": 1.0},
			want: true,
		},
		{
			name: "map missing key",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"b": 1.0},
			want: false,
		},
		{
			name: "different map length",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"a": 1.0, "b": 2.0},
			want: false,
		},
		{
			name: "different map",
			a:    map[string]any{"a": 1.0},
			b:    map[string]any{"a": 2.0},
			want: false,
		},
		{
			name: "same slice",
			a:    []any{1.0, 2.0},
			b:    []any{1.0, 2.0},
			want: true,
		},
		{
			name: "same nested slice",
			a:    []any{[]any{1.0}, []any{2.0}},
			b:    []any{[]any{1.0}, []any{2.0}},
			want: true,
		},
		{
			name: "same slice of maps",
			a:    []any{map[string]any{"a": 1.0}},
			b:    []any{map[string]any{"a": 1.0}},
			want: true,
		},
		{
			name: "different slice length",
			a:    []any{1.0, 2.0},
			b:    []any{1.0},
			want: false,
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
	assert.Equal(t, "arr", typeVar([]any{}))
	assert.Equal(t, unknownType, typeVar(struct{}{}))

	assert.Equal(t, float64(1), normalizeValue(1))
	assert.Equal(t, float64(1), normalizeValue(1.0))
}

func TestBuildDiff(t *testing.T) {
	assert.Equal(t, diff.Node{
		TypeDiff: diff.Added,
		OldValue: nil,
		NewValue: 1.0,
	}, BuildDiff(diff.Added, 1, 1.0))

	assert.Equal(t, diff.Node{
		TypeDiff: diff.Changed,
		OldValue: 1.0,
		NewValue: 2.0,
	}, BuildDiff(diff.Changed, 1, 2))
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
			want: diff.Node{
				TypeDiff: diff.Changed,
				OldValue: true,
				NewValue: false,
			},
		},
		{
			name: "unchanged nested map",
			a: map[string]any{
				"k": "v",
			},
			b: map[string]any{
				"k": "v",
			},
			want: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "k",
						TypeDiff: diff.Unchanged,
						OldValue: "v",
					},
				},
			},
		},
		{
			name: "changed nested slice",
			a:    []any{1.0, 2.0},
			b:    []any{1.0, 3.0},
			want: diff.Node{
				TypeDiff: diff.Changed,
				OldValue: []any{1.0, 2.0},
				NewValue: []any{1.0, 3.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RecursiveGendiff(tt.a, tt.b))
		})
	}
}
