package formatters_test

import (
	// "fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code/diff"
	"code/formatters"
)

func TestFormatPlainNoTrailingNewline(t *testing.T) {
	diffNode := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "key1",
				TypeDiff: diff.Added,
				NewValue: "added value",
			},
		},
	}

	got, err := formatters.PrintDiff(diffNode, formatters.Plain)
	require.NoError(t, err, "PrintDiff() returned an error: %v", err)
	require.Equal(t, "Property 'key1' was added with value: 'added value'", got)
}

func TestFormatPlain(t *testing.T) {
	tests := []struct {
		name string
		diff diff.Node
		want string
		path string
	}{
		{
			name: "plain format flat",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key1",
						TypeDiff: diff.Unchanged,
						OldValue: "value",
					},
					{
						Key:      "key2",
						TypeDiff: diff.Unchanged,
						OldValue: 42.0,
					},
					{
						Key:      "key3",
						TypeDiff: diff.Unchanged,
						OldValue: nil,
					},
					{
						Key:      "key4",
						TypeDiff: diff.Unchanged,
						OldValue: true,
					},
					{
						Key:      "key5",
						TypeDiff: diff.Added,
						NewValue: "added value",
					},
					{
						Key:      "key6",
						TypeDiff: diff.Removed,
						NewValue: "removed value",
					},
					{
						Key:      "key7",
						TypeDiff: diff.Changed,
						OldValue: "old value",
						NewValue: "new value",
					},
				},
			},
			path: filepath.Join("testdata", "fixture", "plainFlat.txt"),
		},
		{
			name: "plain format nested",
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
							{
								Key:      "nestedKey2",
								TypeDiff: diff.Added,
								NewValue: "addedNestedValue",
							},
							{
								Key:      "nestedKey3",
								TypeDiff: diff.Removed,
								OldValue: "removedNestedValue",
							},
							{
								Key:      "nestedKey4",
								TypeDiff: diff.Nested,
								Children: []diff.Node{
									{
										Key:      "deepNestedKey",
										TypeDiff: diff.Changed,
										OldValue: "oldDeepNestedValue",
										NewValue: "newDeepNestedValue",
									},
								},
							},
						},
					},
				},
			},
			path: filepath.Join("testdata", "fixture", "plainNested.txt"),
		},
		{
			name: "plain format empty nested object",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{},
			},
			want: "",
		},
		{
			name: "plain format added complex value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "complexKey",
						TypeDiff: diff.Added,
						NewValue: []any{1.0, 2.0},
					},
				},
			},
			want: "Property 'complexKey' was added with value: [complex value]",
		},
		{
			name: "plain format added nil value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "nullKey",
						TypeDiff: diff.Added,
						NewValue: nil,
					},
				},
			},
			want: "Property 'nullKey' was added with value: null",
		},
		{
			name: "plain format added bool value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "boolKey",
						TypeDiff: diff.Added,
						NewValue: true,
					},
				},
			},
			want: "Property 'boolKey' was added with value: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(tt.diff, formatters.Plain)
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
