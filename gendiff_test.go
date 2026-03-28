package code_test

import (
	"code"
	"code/parser"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}

func parseFile(t *testing.T, path string) any {
	t.Helper()

	data, err := parser.Parse(path)
	require.NoError(t, err, "failed to parse file %s", path)

	return data
}

func TestGendiff(t *testing.T) {
	tests := []struct {
		name         string
		pathA        string
		pathB        string
		expectedPath string
	}{
		{
			name:         "AB",
			pathA:        "testdata/fixture/fileA.json",
			pathB:        "testdata/fixture/fileB.json",
			expectedPath: "testdata/fixture/expectedAB.txt",
		},
		{
			name:         "AC",
			pathA:        "testdata/fixture/fileA.json",
			pathB:        "testdata/fixture/fileC.json",
			expectedPath: "testdata/fixture/expectedAC.txt",
		},
		{
			name:         "CD",
			pathA:        "testdata/fixture/fileC.json",
			pathB:        "testdata/fixture/fileD.json",
			expectedPath: "testdata/fixture/expectedCD.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataA := parseFile(t, tt.pathA)
			dataB := parseFile(t, tt.pathB)

			diff := strings.TrimSpace(code.GenDiff(dataA, dataB))
			expected := readFileToString(t, tt.expectedPath)
			assert.Equal(t, expected, diff, "diff mismatch for case %s", tt.name)
		})
	}
}
