package formatters_test

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"code/diff"
	"code/formatters"
)

// readFileToString читает файл и возвращает обрезанную строку.
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

// runDiffTests выполняет общий цикл тестирования для заданного формата.
func runDiffTests(t *testing.T, format formatters.PrintFormat, tests []diffTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(tt.diff, format)
			require.NoError(t, err, "PrintDiff() returned an error: %v", err)
			actual := strings.TrimSpace(got)

			if tt.path != "" {
				expected := readFileToString(t, tt.path)
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Errorf("PrintDiff() mismatch (-want +got):\n%s", diff)
				}
			} else if actual != tt.want {
				t.Errorf("PrintDiff() = %v, want %v", actual, tt.want)
			}
		})
	}
}
