package parsers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/diff"
)

var fixtureDir = filepath.Join("..", "testdata", "fixture")

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected diff.Value
	}{
		{
			name: "valid json",
			path: filepath.Join(fixtureDir, "flat", "fileA.json"),
			expected: diff.Map{
				"host":    diff.String("hexlet.io"),
				"timeout": diff.Number(50),
				"proxy":   diff.String("123.234.53.22"),
				"follow":  diff.Boolean(false),
			},
		},
		{
			name: "valid yml",
			path: filepath.Join(fixtureDir, "flat", "fileA.yml"),
			expected: diff.Map{
				"host":    diff.String("hexlet.io"),
				"timeout": diff.Number(50),
				"proxy":   diff.String("123.234.53.22"),
				"follow":  diff.Boolean(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseFileNotFound(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "non existent file", path: "non-existent.json"},
		{name: "directory", path: "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "parse")
		})
	}
}

func TestParseInvalidContent(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "invalid json", path: filepath.Join(fixtureDir, "invalid", "invalid.json")},
		{name: "invalid yml", path: filepath.Join(fixtureDir, "invalid", "invalid.yaml")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "parse")
			assert.Contains(t, err.Error(), filepath.Base(tt.path))
		})
	}
}

func TestParseUnsupportedFileType(t *testing.T) {
	_, err := Parse(filepath.Join(fixtureDir, "invalid", "fileA.xml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse")
	assert.Contains(t, err.Error(), "fileA.xml")
	assert.True(t, errors.Is(err, ErrUnsupportedFileType))
}

func TestParseMissingExtension(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("failed to remove temp directory: %v", err)
		}
	}()

	path := filepath.Join(tempDir, "noext")
	err = os.WriteFile(path, []byte("{}"), 0o644)
	require.NoError(t, err)

	_, err = Parse(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing file extension")
}

func TestParsePermissionDenied(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("failed to remove temp directory: %v", err)
		}
	}()

	path := filepath.Join(tempDir, "fileZ.json")

	err = os.WriteFile(path, []byte("{}"), 0o644)
	require.NoError(t, err)

	if err := os.Chmod(path, 0); err != nil {
		t.Skip("chmod not supported on this platform")
	}

	_, err = Parse(path)
	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrPermission))
}
