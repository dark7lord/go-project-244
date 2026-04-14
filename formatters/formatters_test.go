package formatters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code"
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
				"key6": code.Diff{
					TypeDiff: code.DiffTypeAdded,
					NewValue: "added value",
				},
				"key7": code.Diff{
					TypeDiff: code.DiffTypeRemoved,
					OldValue: "removed value",
				},
				"key8": code.Diff{
					TypeDiff: code.DiffTypeUnchanged,
					OldValue: "unchanged value",
					NewValue: "unchanged value",
				},
				"key9": code.Diff{
					TypeDiff: code.DiffTypeChanged,
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
				"key3": code.Diff{
					TypeDiff: code.DiffTypeAdded,
					NewValue: map[string]any{
						"addedNestedKey": "addedNestedValue",
					},
				},
				"key4": code.Diff{
					TypeDiff: code.DiffTypeRemoved,
					OldValue: map[string]any{
						"removedNestedKey": "removedNestedValue",
					},
				},
				"key5": code.Diff{
					TypeDiff: code.DiffTypeChanged,
					OldValue: map[string]any{
						"oldNestedKey": "oldNestedValue",
					},
					NewValue: map[string]any{
						"newNestedKey": "newNestedValue",
					},
				},
				"key6": code.Diff{
					TypeDiff: code.DiffTypeUnchanged,
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
				"key5": code.Diff{
					TypeDiff: code.DiffTypeAdded,
					NewValue: "added value",
				},
				"key6": code.Diff{
					TypeDiff: code.DiffTypeRemoved,
					OldValue: "removed value",
				},
				"key7": code.Diff{
					TypeDiff: code.DiffTypeChanged,
					OldValue: "old value",
					NewValue: "new value",
				},
				"key8": code.Diff{
					TypeDiff: code.DiffTypeUnchanged,
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
					"nestedKey2": code.Diff{
						TypeDiff: code.DiffTypeAdded,
						NewValue: "addedNestedValue",
					},
					"nestedKey3": code.Diff{
						TypeDiff: code.DiffTypeRemoved,
						OldValue: "removedNestedValue",
					},
					"nestedKey4": map[string]any{
						"deepNestedKey": code.Diff{
							TypeDiff: code.DiffTypeChanged,
							OldValue: "oldDeepNestedValue",
							NewValue: "newDeepNestedValue",
						},
					},
					"nestedKey5": code.Diff{
						TypeDiff: code.DiffTypeUnchanged,
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
					"nestedKey2": code.Diff{
						TypeDiff: code.DiffTypeAdded,
						NewValue: "added nested value",
					},
				},
				"key6": code.Diff{
					TypeDiff: code.DiffTypeRemoved,
					OldValue: "deleted value",
				},
				"key7": code.Diff{
					TypeDiff: code.DiffTypeChanged,
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
