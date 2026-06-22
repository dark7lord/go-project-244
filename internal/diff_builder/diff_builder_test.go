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

func TestGenMapDiff(t *testing.T) {
	tests := []struct {
		name string
		mapA diff.Map
		mapB diff.Map
		want diff.Node
	}{
		{
			name: "flat add",
			mapA: diff.Map{},
			mapB: diff.Map{testKey: diff.Number(42)},
			want: nestedDiff(
				diff.Node{Key: testKey, TypeDiff: diff.Added, NewValue: 42.0},
			),
		},
		{
			name: "flat remove",
			mapA: diff.Map{testKey: diff.String("value")},
			mapB: diff.Map{},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Removed, OldValue: "value"},
			),
		},
		{
			name: "flat changed",
			mapA: diff.Map{testKey: diff.Boolean(true)},
			mapB: diff.Map{testKey: diff.Boolean(false)},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Changed, OldValue: true, NewValue: false},
			),
		},
		{
			name: "flat unchanged",
			mapA: diff.Map{testKey: diff.Null{}},
			mapB: diff.Map{testKey: diff.Null{}},
			want: nestedDiff(
				diff.Node{
					Key: testKey, TypeDiff: diff.Unchanged, OldValue: nil},
			),
		},
		{
			name: "nested object",
			mapA: diff.Map{testNested: diff.Map{"a": diff.Number(1)}},
			mapB: diff.Map{testNested: diff.Map{"a": diff.Number(1), "b": diff.Number(2)}},
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
			mapA: diff.Map{testArrayKey: diff.Slice{diff.Number(1), diff.Number(2)}},
			mapB: diff.Map{testArrayKey: diff.Slice{diff.Number(1), diff.Number(3), diff.Number(4)}},
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
			got := genMapDiff(tt.mapA, tt.mapB)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name string
		a    diff.Value
		b    diff.Value
		want bool
	}{
		{
			name: "same bool",
			a:    diff.Boolean(true),
			b:    diff.Boolean(true),
			want: true,
		},
		{
			name: "different bool",
			a:    diff.Boolean(true),
			b:    diff.Boolean(false),
			want: false,
		},
		{
			name: "different type",
			a:    diff.Number(1),
			b:    diff.String("1"), want: false},
		{
			name: "same map",
			a:    diff.Map{"a": diff.Number(1)},
			b:    diff.Map{"a": diff.Number(1)},
			want: true,
		},
		{
			name: "different map",
			a:    diff.Map{"a": diff.Number(1)},
			b:    diff.Map{"a": diff.Number(2)},
			want: false,
		},
		{
			name: "maps same length different keys",
			a:    diff.Map{"a": diff.Number(1), "b": diff.Number(2)},
			b:    diff.Map{"a": diff.Number(1), "c": diff.Number(3)},
			want: false,
		},
		{
			name: "same slice",
			a:    diff.Slice{diff.Number(1), diff.Number(2)},
			b:    diff.Slice{diff.Number(1), diff.Number(2)},
			want: true,
		},
		{
			name: "different slice",
			a:    diff.Slice{diff.Number(1), diff.Number(2)},
			b:    diff.Slice{diff.Number(1), diff.Number(3)},
			want: false,
		},
		{
			name: "unknown type",
			a:    diff.Null{},
			b:    diff.Null{},
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
	v := diff.UnknownValue()
	assert.False(t, isEqual(v, v))
}

func TestBuildDiff(t *testing.T) {
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Added, OldValue: nil, NewValue: 1.0},
		buildDiff(diff.Added, diff.Number(1), diff.Number(1)),
	)
	assert.Equal(t,
		diff.Node{TypeDiff: diff.Changed, OldValue: 1.0, NewValue: 2.0},
		buildDiff(diff.Changed, diff.Number(1), diff.Number(2)),
	)
}

func TestRecursiveGendiff(t *testing.T) {
	tests := []struct {
		name string
		a    diff.Value
		b    diff.Value
		want diff.Node
	}{
		{
			name: "changed primitive",
			a:    diff.Boolean(true),
			b:    diff.Boolean(false),
			want: diff.Node{TypeDiff: diff.Changed, OldValue: true, NewValue: false},
		},
		{
			name: "unchanged nested map",
			a:    diff.Map{"k": diff.String("v")},
			b:    diff.Map{"k": diff.String("v")},
			want: nestedDiff(
				diff.Node{Key: "k", TypeDiff: diff.Unchanged, OldValue: "v"},
			),
		},
		{
			name: "changed nested slice",
			a:    diff.Slice{diff.Number(1), diff.Number(2)},
			b:    diff.Slice{diff.Number(1), diff.Number(3)},
			want: diff.Node{TypeDiff: diff.Changed, OldValue: []any{1.0, 2.0}, NewValue: []any{1.0, 3.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RecursiveGendiff(tt.a, tt.b))
		})
	}
}
