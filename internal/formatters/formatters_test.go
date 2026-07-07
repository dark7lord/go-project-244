package formatters_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestFormat(t *testing.T) {
	singleAddedNodeExpected := readFileToString(t, filepath.Join("testdata", "fixtures", "singleAddedNode.json"))

	got, err := formatters.Format(string(formatters.JSON), singleAddedDiff())
	require.NoError(t, err)
	requireOutputEqual(t, formatters.JSON, singleAddedNodeExpected, got)
}

func TestFormatUnsupported(t *testing.T) {
	_, err := formatters.Format("xml", diff.Node{})
	require.Error(t, err)
}
