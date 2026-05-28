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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(tt.diff, formatters.Plain)
			require.NoError(t, err, "PrintDiff() returned an error: %v", err)
			actual := strings.TrimSpace(got)

			// fmt.Println(tt.diff)
			// fmt.Println(actual)

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
