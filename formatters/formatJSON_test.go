package formatters_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/diff"
	"code/formatters"
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

func TestFormatJSONMarshalError(t *testing.T) {
	_, err := formatters.PrintDiff(diff.Node{
		TypeDiff: diff.Unchanged,
		OldValue: make(chan int),
	}, formatters.JSON)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported type")
}
