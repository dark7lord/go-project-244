package formatters_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code/diff"
	"code/formatters"
)

func TestFormatStylish(t *testing.T) {
	tests := []struct {
		name string
		diff diff.Node
		want string
		path string
	}{
		{
			name: "primitive",
			diff: diff.Node{
				TypeDiff: diff.Unchanged,
				OldValue: "value",
			},
			want: "value",
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
			name: "flat map",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key1",
						TypeDiff: diff.Unchanged,
						OldValue: "string value",
					},
					{
						Key:      "key2",
						TypeDiff: diff.Unchanged,
						OldValue: 42.0,
					},
					{
						Key:      "key3",
						TypeDiff: diff.Unchanged,
						OldValue: "",
					},
					{
						Key:      "key4",
						TypeDiff: diff.Unchanged,
						OldValue: nil,
					},
					{
						Key:      "key5",
						TypeDiff: diff.Unchanged,
						OldValue: true,
					},
					{
						Key:      "key6",
						TypeDiff: diff.Added,
						NewValue: "added value",
					},
					{
						Key:      "key7",
						TypeDiff: diff.Removed,
						OldValue: "removed value",
					},
					{
						Key:      "key8",
						TypeDiff: diff.Unchanged,
						OldValue: "unchanged value",
					},
					{
						Key:      "key9",
						TypeDiff: diff.Changed,
						OldValue: "old value",
						NewValue: "new value",
					},
				},
			},
			path: filepath.Join("testdata", "fixture", "flatMap.diff"),
		},
		{
			name: "nested map",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key1",
						TypeDiff: diff.Unchanged,
						OldValue: "value1",
					},
					{
						Key:      "key2",
						TypeDiff: diff.Nested,
						Children: []diff.Node{
							{
								Key:      "nestedKey",
								TypeDiff: diff.Unchanged,
								OldValue: "nestedValue",
							},
						},
					},
					{
						Key:      "key3",
						TypeDiff: diff.Added,
						NewValue: map[string]any{
							"addedNestedKey": "addedNestedValue",
						},
					},
					{
						Key:      "key4",
						TypeDiff: diff.Removed,
						OldValue: map[string]any{
							"removedNestedKey": "removedNestedValue",
						},
					},
					{
						Key:      "key5",
						TypeDiff: "changed",
						OldValue: map[string]any{
							"oldNestedKey": "oldNestedValue",
						},
						NewValue: map[string]any{
							"newNestedKey": "newNestedValue",
						},
					},
					{
						Key:      "key6",
						TypeDiff: diff.Nested,
						Children: []diff.Node{
							{
								Key:      "unchangedNestedKey",
								TypeDiff: diff.Unchanged,
								OldValue: "unchangedNestedValue",
							},
						},
					},
				},
			},
			path: filepath.Join("testdata", "fixture", "nestedMap.diff"),
		},
		{
			name: "empty nested object",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{},
			},
			want: "{}",
		},
		{
			name: "added empty slice",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "arr",
						TypeDiff: diff.Added,
						NewValue: []any{},
					},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(tt.diff, formatters.Stylish)
			require.NoError(t, err, "PrintDiff() returned an error: %v", err)
			actual := strings.TrimSpace(got)

			if tt.path != "" {
				expected := readFileToString(t, tt.path)
				if actual != expected {
					diff := cmp.Diff(expected, actual)
					t.Errorf("PrintDiff() mismatch (-want +got):\n%s", diff)
				}
			} else if actual != tt.want {
				t.Errorf("PrintDiff() = %v, want %v", actual, tt.want)
			}
		})
	}
}
