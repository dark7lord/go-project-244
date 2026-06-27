package formatters_test

import (
	"path/filepath"
	"testing"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestStylishFormatter(t *testing.T) {
	tests := []diffTestCase{
		{
			name: "primitive",
			diff: diff.Node{
				TypeDiff: diff.Unchanged,
				OldValue: testNamePrimitiveValue,
			},
			want: testNamePrimitiveValue,
		},
		{
			name: "null",
			diff: diff.Node{
				TypeDiff: diff.Unchanged,
				OldValue: nil,
			},
			want: "null",
		},
		{
			name: "flat all types",
			diff: flatDiffAll,
			path: filepath.Join("testdata", "fixtures", "flatMap.diff"),
		},
		{
			name: testNameNestedAll,
			diff: nestedDiffAll,
			path: filepath.Join("testdata", "fixtures", "nestedMap.diff"),
		},
		{
			name: testNameEmptyNested,
			diff: emptyNestedDiff,
			want: "{}",
		},
		{
			name: "added empty slice",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{Key: "arr", TypeDiff: diff.Added, NewValue: []any{}},
				},
			},
			want: "{\n  + arr: []\n}",
		},
		{
			name: "added empty map",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "obj",
						TypeDiff: diff.Added,
						NewValue: map[string]any{},
					},
				},
			},
			want: "{\n  + obj: {}\n}",
		},
		{
			name: "added non-empty slice",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "arr",
						TypeDiff: diff.Added,
						NewValue: []any{"x"},
					},
				},
			},
			want: "{\n  + arr: [\n        x\n    ]\n}",
		},
	}
	runDiffTests(t, formatters.Stylish, tests)
}
