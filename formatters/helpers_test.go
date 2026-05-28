package formatters_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}
