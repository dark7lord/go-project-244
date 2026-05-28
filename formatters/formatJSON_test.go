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

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name string
		diff diff.Node
		want string
		path string
	}{
		{
			name: "json format",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "key1",
						TypeDiff: diff.Unchanged,
						OldValue: "string value 1",
					},
					{
						Key:      "key2",
						TypeDiff: diff.Unchanged,
						OldValue: 42.0,
					},
					{
						Key:      "key3",
						TypeDiff: diff.Unchanged,
						OldValue: false,
					},
					{
						Key:      "key4",
						TypeDiff: diff.Unchanged,
						OldValue: nil,
					},
					{
						Key:      "key5",
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
								NewValue: "added nested value",
							},
						},
					},
					{
						Key:      "key6",
						TypeDiff: diff.Removed,
						OldValue: "deleted value",
					},
					{
						Key:      "key7",
						TypeDiff: diff.Changed,
						OldValue: "old value",
						NewValue: "new value",
					},
				},
			},
			path: filepath.Join("testdata", "fixture", "jsonNested.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(tt.diff, formatters.JSON)
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
