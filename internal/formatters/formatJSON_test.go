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
			want: `{}`,
		},
	}
	runDiffTests(t, formatters.JSON, tests)
}

func TestFormatJSONNoError(t *testing.T) {
	_, err := formatters.PrintDiff(diff.Node{
		TypeDiff: diff.Unchanged,
		OldValue: "test",
	}, formatters.JSON)

	require.NoError(t, err)
}

func TestFormatJSONUnknownChildKind(t *testing.T) {
	node := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{Key: "unknown", TypeDiff: diff.Kind("bogus")},
		},
	}
	result, err := formatters.PrintDiff(node, formatters.JSON)
	require.NoError(t, err)
	require.Equal(t, "{}", result)
}

func TestFormatJSONUnknownKind(t *testing.T) {
	result, err := formatters.PrintDiff(diff.Node{TypeDiff: diff.Kind("bogus")}, formatters.JSON)
	require.NoError(t, err)
	require.Equal(t, "{}", result)
}
