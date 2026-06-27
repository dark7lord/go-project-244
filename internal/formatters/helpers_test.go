package formatters_test

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return strings.TrimSpace(string(data))
}

type diffTestCase struct {
	name string
	diff diff.Node
	want string
	path string
}

func requireOutputEqual(t *testing.T, format formatters.OutputFormat, expected, actual string) {
	t.Helper()

	if format == formatters.JSON {
		require.JSONEq(t, expected, actual)
		return
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("formatter output mismatch (-want +got):\n%s", diff)
	}
}

func runDiffTests(t *testing.T, format formatters.OutputFormat, tests []diffTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.Format(string(format), tt.diff)
			require.NoError(t, err, "Format() returned an error: %v", err)
			actual := strings.TrimSpace(got)

			if tt.path != "" {
				expected := readFileToString(t, tt.path)
				requireOutputEqual(t, format, expected, actual)
			} else {
				requireOutputEqual(t, format, tt.want, actual)
			}
		})
	}
}
