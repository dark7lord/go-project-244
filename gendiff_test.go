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
	expectedGGDiff = "expected/expectedGG.diff"
	expectedGHDiff = "expected/expectedGH.diff"
)

func fixturePath(name string) string {
	return filepath.Join("testdata", "gendiff", "fixture", name)
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
		{
			name:         "GG same files",
			pathA:        fixturePath(fixtureSameAJSON),
			pathB:        fixturePath(fixtureSameAJSON),
			expectedPath: fixturePath(expectedGGDiff),
		},
		{
			name:         "GH identical files",
			pathA:        fixturePath(fixtureSameAJSON),
			pathB:        fixturePath(fixtureSameBJSON),
			expectedPath: fixturePath(expectedGHDiff),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenDiff(tt.pathA, tt.pathB, "stylish")
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
	}{
		{
			name:        "invalid first file",
			pathA:       fixturePath(fixtureInvalidJSON),
			pathB:       fixturePath(fixtureFlatAJSON),
			format:      "stylish",
			expectedErr: "",
		},
		{
			name:        "invalid second file",
			pathA:       fixturePath(fixtureFlatAJSON),
			pathB:       fixturePath(fixtureInvalidYAML),
			format:      "stylish",
			expectedErr: "",
		}, {
			name:        "unsupported format",
			pathA:       fixturePath(fixtureFlatAJSON),
			pathB:       fixturePath(fixtureFlatBJSON),
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
