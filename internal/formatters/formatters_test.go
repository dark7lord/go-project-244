package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestFormat(t *testing.T) {
	got, err := formatters.Format(string(formatters.JSON), singleAddedDiff())
	require.NoError(t, err)
	requireOutputEqual(t, formatters.JSON, singleAddedNodeJSON, got)
}

func TestFormatUnsupported(t *testing.T) {
	_, err := formatters.Format("xml", diff.Node{})
	require.Error(t, err)
}
