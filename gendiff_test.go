package code

import (
	"os"
	"path/filepath"
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

func fixturePath(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func TestGendiff(t *testing.T) {
	tests := []struct {
		name         string
		pathA        string
		pathB        string
		expectedPath string
	}{
		{
			name:         "AB flat json",
			pathA:        fixturePath("fileA.json"),
			pathB:        fixturePath("fileB.json"),
			expectedPath: fixturePath("expectedAB.diff"),
		},
		{
			name:         "AC",
			pathA:        fixturePath("fileA.json"),
			pathB:        fixturePath("fileC.json"),
			expectedPath: fixturePath("expectedAC.diff"),
		},
		{
			name:         "CD",
			pathA:        fixturePath("fileC.json"),
			pathB:        fixturePath("fileD.json"),
			expectedPath: fixturePath("expectedCD.diff"),
		},
		{
			name:         "AA json vs yaml",
			pathA:        fixturePath("fileA.json"),
			pathB:        fixturePath("fileA.yml"),
			expectedPath: fixturePath("expectedAA.diff"),
		},
		{
			name:         "AA yaml vs json",
			pathA:        fixturePath("fileA.yml"),
			pathB:        fixturePath("fileA.json"),
			expectedPath: fixturePath("expectedAA.diff"),
		},
		{
			name:         "AB yaml",
			pathA:        fixturePath("fileA.yml"),
			pathB:        fixturePath("fileB.yml"),
			expectedPath: fixturePath("expectedAB.diff"),
		},
		{
			name:         "nested json files",
			pathA:        fixturePath("fileE.json"),
			pathB:        fixturePath("fileF.json"),
			expectedPath: fixturePath("expectedEF.diff"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := GenDiff(tt.pathA, tt.pathB, "stylish")
			actual := strings.TrimSpace(result)
			expected := readFileToString(t, tt.expectedPath)

			assert.Equal(t, expected, actual, "diff mismatch for case %s", tt.name)
		})
	}
}

func TestGenDiffErrors(t *testing.T) {
	tests := []struct {
		name        string
		pathA       string
		pathB       string
		format      string
		expectedErr string
	}{
		{
			name:        "invalid first file",
			pathA:       fixturePath("invalid.json"),
			pathB:       fixturePath("fileA.json"),
			format:      "stylish",
			expectedErr: "",
		},
		{
			name:        "invalid second file",
			pathA:       fixturePath("fileA.json"),
			pathB:       fixturePath("invalid.yaml"),
			format:      "stylish",
			expectedErr: "",
		},
		{
			name:        "unsupported format",
			pathA:       fixturePath("fileA.json"),
			pathB:       fixturePath("fileB.json"),
			format:      "xml",
			expectedErr: "unsupported format: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenDiff(tt.pathA, tt.pathB, tt.format)

			if tt.expectedErr == "" {
				require.Error(t, err)
				assert.Empty(t, result)
				return
			}

			require.EqualError(t, err, tt.expectedErr)
			assert.Empty(t, result)
		})
	}
}
