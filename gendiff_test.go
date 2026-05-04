package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/diff"
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

func TestGenSliceDiff(t *testing.T) {
	tests := []struct {
		name  string
		left  []any
		right []any
		want  []any
	}{
		{
			name:  "same slices",
			left:  []any{1, 2, 3},
			right: []any{1, 2, 3},
			want: []any{
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 1.0,
				},
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 2.0,
				},
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 3.0,
				},
			},
		},
		{
			name:  "flat arrays with adding",
			left:  []any{1, 2, 3},
			right: []any{1, 4, 3, 5},
			want: []any{
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 1.0,
				},
				diff.Diff{
					TypeDiff: diff.Changed,
					OldValue: 2.0,
					NewValue: 4.0,
				},
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 3.0,
				},
				diff.Diff{
					TypeDiff: diff.Added,
					NewValue: 5.0,
				},
			},
		},
		{
			name:  "flat arrays with removing",
			left:  []any{1, 2, 3},
			right: []any{1, 2},
			want: []any{
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 1.0,
				},
				diff.Diff{
					TypeDiff: diff.Unchanged,
					OldValue: 2.0,
				},
				diff.Diff{
					TypeDiff: diff.Removed,
					OldValue: 3.0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genSliceDiff(tt.left, tt.right)
			assert.Equal(t, got, tt.want, "got: \n\t%v\n, but want: \n\t%v\n", got, tt.want)
		})
	}
}
