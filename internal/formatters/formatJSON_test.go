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
			path: filepath.Join("testdata", "fixture", "jsonNested.json"),
		},
		{
			name: testNameEmptyNested,
			diff: emptyNestedDiff,
			want: "{}",
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
