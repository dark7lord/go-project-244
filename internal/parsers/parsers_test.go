package parsers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtureDir = filepath.Join("..", "..", "testdata", "fixtures")

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{
			name: "valid json",
			path: filepath.Join(fixtureDir, "flat", "fileA.json"),
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": float64(50),
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name: "valid yml",
			path: filepath.Join(fixtureDir, "flat", "fileA.yml"),
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": float64(50),
				"proxy":   "123.234.53.22",
				"follow":  false,
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

func TestNormalizeYAML(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want any
	}{
		{name: "int", v: int(42), want: float64(42)},
		{name: "string", v: "hello", want: "hello"},
		{name: "bool", v: true, want: true},
		{name: "float64", v: 3.14, want: 3.14},
		{name: "nil", v: nil, want: nil},
		{name: "slice",
			v:    []any{int(1), "x", true, nil},
			want: []any{float64(1), "x", true, nil}},
		{name: "nested map",
			v:    map[string]any{"a": int(1), "b": map[string]any{"c": int(2)}},
			want: map[string]any{"a": float64(1), "b": map[string]any{"c": float64(2)}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeYAML(tt.v))
		})
	}
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

func TestParseRejectsRootArray(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("failed to remove temp directory: %v", err)
		}
	}()

	cases := []struct {
		name string
		path string
		data []byte
	}{
		{name: "json array", path: filepath.Join(tempDir, "array.json"), data: []byte(`["a", "b"]`)},
		{name: "yaml array", path: filepath.Join(tempDir, "array.yaml"), data: []byte("- a\n- b\n")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tt.path, tt.data, 0o644))

			_, err := Parse(tt.path)
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrUnsupportedRootType) || assert.ErrorContains(t, err, "cannot unmarshal"))
		})
	}
}
