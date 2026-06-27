package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/formatters"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}

const (
	fixtureFlatAJSON   = "flat/fileA.json"
	fixtureFlatBJSON   = "flat/fileB.json"
	fixtureFlatCJSON   = "flat/fileC.json"
	fixtureFlatDJSON   = "flat/fileD.json"
	fixtureNestedAJSON = "nested/fileE.json"
	fixtureNestedBJSON = "nested/fileF.json"
	fixtureSameAJSON   = "flat/fileG.json"
	fixtureSameBJSON   = "flat/fileH.json"
	fixtureAAYML       = "flat/fileA.yml"
	fixtureBYML        = "flat/fileB.yml"
	fixtureInvalidJSON = "invalid/invalid.json"
	fixtureInvalidYAML = "invalid/invalid.yaml"

	expectedABDiff = "expected/expectedAB.diff"
	expectedACDiff = "expected/expectedAC.diff"
	expectedCDDiff = "expected/expectedCD.diff"
	expectedAADiff = "expected/expectedAA.diff"
	expectedEFDiff = "expected/expectedEF.diff"
)

func fixturePath(name string) string {
	return filepath.Join("testdata", "fixtures", name)
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
			pathA:        fixturePath(fixtureFlatAJSON),
			pathB:        fixturePath(fixtureFlatBJSON),
			expectedPath: fixturePath(expectedABDiff),
		},
		{
			name:         "AC",
			pathA:        fixturePath(fixtureFlatAJSON),
			pathB:        fixturePath(fixtureFlatCJSON),
			expectedPath: fixturePath(expectedACDiff),
		},
		{
			name:         "CD",
			pathA:        fixturePath(fixtureFlatCJSON),
			pathB:        fixturePath(fixtureFlatDJSON),
			expectedPath: fixturePath(expectedCDDiff),
		},
		{
			name:         "AA json vs yaml",
			pathA:        fixturePath(fixtureFlatAJSON),
			pathB:        fixturePath(fixtureAAYML),
			expectedPath: fixturePath(expectedAADiff),
		},
		{
			name:         "AA yaml vs json",
			pathA:        fixturePath(fixtureAAYML),
			pathB:        fixturePath(fixtureFlatAJSON),
			expectedPath: fixturePath(expectedAADiff),
		},
		{
			name:         "AB yaml",
			pathA:        fixturePath(fixtureAAYML),
			pathB:        fixturePath(fixtureBYML),
			expectedPath: fixturePath(expectedABDiff),
		},
		{
			name:         "nested json files",
			pathA:        fixturePath(fixtureNestedAJSON),
			pathB:        fixturePath(fixtureNestedBJSON),
			expectedPath: fixturePath(expectedEFDiff),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenDiff(tt.pathA, tt.pathB, string(formatters.Stylish))
			require.NoError(t, err)
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
		contains    bool
	}{
		{
			name:        "invalid first file",
			pathA:       fixturePath(fixtureInvalidJSON),
			pathB:       fixturePath(fixtureFlatAJSON),
			format:      string(formatters.Stylish),
			expectedErr: "",
		},
		{
			name:        "invalid second file",
			pathA:       fixturePath(fixtureFlatAJSON),
			pathB:       fixturePath(fixtureInvalidYAML),
			format:      string(formatters.Stylish),
			expectedErr: "",
		},
		{
			name:        "unsupported format",
			pathA:       fixturePath(fixtureFlatAJSON),
			pathB:       fixturePath(fixtureFlatBJSON),
			format:      "xml",
			expectedErr: "unsupported format: xml",
		},
		{
			name:        "array root in first file",
			pathA:       fixturePath(fixtureSameAJSON),
			pathB:       fixturePath(fixtureFlatAJSON),
			format:      string(formatters.Stylish),
			expectedErr: "cannot unmarshal array into Go value of type map[string]interface {}",
			contains:    true,
		},
		{
			name:        "array root in second file",
			pathA:       fixturePath(fixtureFlatAJSON),
			pathB:       fixturePath(fixtureSameBJSON),
			format:      string(formatters.Stylish),
			expectedErr: "cannot unmarshal array into Go value of type map[string]interface {}",
			contains:    true,
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

			if tt.contains {
				require.ErrorContains(t, err, tt.expectedErr)
				assert.Empty(t, result)
				return
			}

			require.EqualError(t, err, tt.expectedErr)
			assert.Empty(t, result)
		})
	}
}
