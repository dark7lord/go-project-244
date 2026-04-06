package code_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code"
	"code/parsers"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}

func parseFile(t *testing.T, path string) any {
	t.Helper()

	data, err := parsers.Parse(path)
	require.NoError(t, err, "failed to parse file %s", path)

	return data
}

const (
	fixtureDir = "testdata/fixture"
	fileA      = fixtureDir + "/fileA.json"
	fileB      = fixtureDir + "/fileB.json"
	fileC      = fixtureDir + "/fileC.json"
	fileD      = fixtureDir + "/fileD.json"
	fileE      = fixtureDir + "/fileE.json"
	fileF      = fixtureDir + "/fileF.json"
	fileAYaml  = fixtureDir + "/fileA.yml"
	fileBYaml  = fixtureDir + "/fileB.yml"
	expectedAB = fixtureDir + "/expectedAB.diff"
	expectedAC = fixtureDir + "/expectedAC.diff"
	expectedCD = fixtureDir + "/expectedCD.diff"
	expectedAA = fixtureDir + "/expectedAA.diff"
	expectedEF = fixtureDir + "/expectedEF.diff"
)

func TestGendiffFlat(t *testing.T) {
	tests := []struct {
		name         string
		pathA        string
		pathB        string
		expectedPath string
	}{
		{
			name:         "AB",
			pathA:        fileA,
			pathB:        fileB,
			expectedPath: expectedAB,
		},
		{
			name:         "AC",
			pathA:        fileA,
			pathB:        fileC,
			expectedPath: expectedAC,
		},
		{
			name:         "CD",
			pathA:        fileC,
			pathB:        fileD,
			expectedPath: expectedCD,
		},
		{
			name:         "AA json vs yaml",
			pathA:        fileA,
			pathB:        fileAYaml,
			expectedPath: expectedAA,
		},
		{
			name:         "AA yaml vs json",
			pathA:        fileAYaml,
			pathB:        fileA,
			expectedPath: expectedAA,
		},
		{
			name:         "AB yaml",
			pathA:        fileAYaml,
			pathB:        fileBYaml,
			expectedPath: expectedAB,
		},
		{
			name:         "nested json files",
			pathA:        fileE,
			pathB:        fileF,
			expectedPath: expectedEF,
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
