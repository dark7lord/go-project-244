package formatters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code/diff"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}

func TestPrintDiff(t *testing.T) {
	tests := []struct {
		name   string
		diff   any
		format string
		want   string
		path   string
	}{
		{
			name:   "primitive",
			diff:   "value",
			format: "stylish",
			want:   "value",
		},
		{
			name:   "null",
			diff:   nil,
			format: "stylish",
			want:   "null",
		},
		{
			name: "flat map",
			diff: map[string]any{
				"key1": "string value",
				"key2": 42,
				"key3": "",
				"key4": nil,
				"key5": true,
				"key6": diff.Diff{
					TypeDiff: "added",
					NewValue: "added value",
				},
				"key7": diff.Diff{
					TypeDiff: "removed",
					OldValue: "removed value",
				},
				"key8": diff.Diff{
					TypeDiff: "unchanged",
					OldValue: "unchanged value",
					NewValue: "unchanged value",
				},
				"key9": diff.Diff{
					TypeDiff: "changed",
					OldValue: "old value",
					NewValue: "new value",
				},
			},
			format: "stylish",
			path:   filepath.Join("testdata", "fixture", "flatMap.diff"),
		},
		{
			name: "nested map",
			diff: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"nestedKey": "nestedValue",
				},
				"key3": diff.Diff{
					TypeDiff: "added",
					NewValue: map[string]any{
						"addedNestedKey": "addedNestedValue",
					},
				},
				"key4": diff.Diff{
					TypeDiff: "removed",
					OldValue: map[string]any{
						"removedNestedKey": "removedNestedValue",
					},
				},
				"key5": diff.Diff{
					TypeDiff: "changed",
					OldValue: map[string]any{
						"oldNestedKey": "oldNestedValue",
					},
					NewValue: map[string]any{
						"newNestedKey": "newNestedValue",
					},
				},
				"key6": diff.Diff{
					TypeDiff: "unchanged",
					OldValue: map[string]any{
						"unchangedNestedKey": "unchangedNestedValue",
					},
					NewValue: map[string]any{
						"unchangedNestedKey": "unchangedNestedValue",
					},
				},
			},
			format: "stylish",
			path:   filepath.Join("testdata", "fixture", "nestedMap.diff"),
		},
		{
			name: "plain format flat",
			diff: map[string]any{
				"key1": "value",
				"key2": 42,
				"key3": nil,
				"key4": true,
				"key5": diff.Diff{
					TypeDiff: "added",
					NewValue: "added value",
				},
				"key6": diff.Diff{
					TypeDiff: "removed",
					OldValue: "removed value",
				},
				"key7": diff.Diff{
					TypeDiff: "changed",
					OldValue: "old value",
					NewValue: "new value",
				},
				"key8": diff.Diff{
					TypeDiff: "unchanged",
					OldValue: "unchanged value",
					NewValue: "unchanged value",
				},
			},
			format: "plain",
			path:   filepath.Join("testdata", "fixture", "plainFlat.txt"),
		},
		{
			name: "plain format nested",
			diff: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"nestedKey": "nestedValue",
					"nestedKey2": diff.Diff{
						TypeDiff: "added",
						NewValue: "addedNestedValue",
					},
					"nestedKey3": diff.Diff{
						TypeDiff: "removed",
						OldValue: "removedNestedValue",
					},
					"nestedKey4": map[string]any{
						"deepNestedKey": diff.Diff{
							TypeDiff: "changed",
							OldValue: "oldDeepNestedValue",
							NewValue: "newDeepNestedValue",
						},
					},
					"nestedKey5": diff.Diff{
						TypeDiff: "unchanged",
						OldValue: "unchangedNestedValue",
					},
				},
			},
			format: "plain",
			path:   filepath.Join("testdata", "fixture", "plainNested.txt"),
		},
		{
			name: "json format",
			diff: map[string]any{
				"key1": "string value 1",
				"key2": 42,
				"key3": false,
				"key4": nil,
				"key5": map[string]any{
					"nestedKey": "nestedValue",
					"nestedKey2": diff.Diff{
						TypeDiff: "added",
						NewValue: "added nested value",
					},
				},
				"key6": diff.Diff{
					TypeDiff: "removed",
					OldValue: "deleted value",
				},
				"key7": diff.Diff{
					TypeDiff: "changed",
					OldValue: "old value",
					NewValue: "new value",
				},
			},
			format: "json",
			path:   filepath.Join("testdata", "fixture", "jsonNested.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrintDiff(tt.diff, tt.format)
			actual := strings.TrimSpace(got)

			if tt.path != "" {
				expected := readFileToString(t, tt.path)
				if actual != expected {
					diff := cmp.Diff(expected, actual)
					t.Errorf("PrintDiff() mismatch (-want +got):\n%s", diff)
					// t.Errorf("PrintDiff() = %v, want %v", actual, expected)
				}
			} else if actual != tt.want {
				t.Errorf("PrintDiff() = %v, want %v", actual, tt.want)
			}
		})
	}
}
