package formatters_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestFormatJSON(t *testing.T) {
	tests := []diffTestCase{
		{
			name: testNameNestedAll,
			diff: nestedDiffAll,
			path: filepath.Join("testdata", "fixtures", "jsonNested.json"),
		},
		{
			name: testNameEmptyNested,
			diff: emptyNestedDiff,
			want: "[]",
		},
		{
			name: "nil values",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "nullAdded",
						TypeDiff: diff.Added,
						NewValue: diff.Null{},
					},
					{
						Key:      "nullChanged",
						TypeDiff: diff.Changed,
						OldValue: diff.String("old"),
						NewValue: diff.Null{},
					},
				},
			},
			want: `[
  {
    "key": "nullAdded",
    "type": "added",
    "value": null
  },
  {
    "key": "nullChanged",
    "type": "changed",
    "oldValue": "old",
    "newValue": null
  }
]`,
		},
	}
	runDiffTests(t, formatters.JSON, tests)
}

func TestFormatJSONNoError(t *testing.T) {
	_, err := formatters.PrintDiff(diff.Node{
		TypeDiff: diff.Unchanged,
		OldValue: diff.String("test"),
	}, formatters.JSON)

	require.NoError(t, err)
}

func TestFormatJSONUnknownKind(t *testing.T) {
	result, err := formatters.PrintDiff(diff.Node{TypeDiff: diff.Kind("bogus")}, formatters.JSON)
	require.NoError(t, err)
	require.Equal(t, "[]", result)
}
