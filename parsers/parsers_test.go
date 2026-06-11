package parsers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtureDir = filepath.Join("..", "testdata", "gendiff", "fixture")

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected any
		wantErr  bool
	}{
		{
			name:    "non existent file",
			path:    "non-existent.json",
			wantErr: true,
		},
		{
			name:    "directory",
			path:    ".",
			wantErr: true,
		},
		{
			name:    "invalid json",
			path:    filepath.Join(fixtureDir, "invalid.json"),
			wantErr: true,
		},
		{
			name:    "invalid yml",
			path:    filepath.Join(fixtureDir, "invalid.yaml"),
			wantErr: true,
		},
		{
			name: "valid json",
			path: filepath.Join(fixtureDir, "flat", "fileA.json"),
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name: "valid yml",
			path: filepath.Join(fixtureDir, "flat", "fileA.yml"),
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name:    "unsupported file type",
			path:    filepath.Join(fixtureDir, "invalid", "fileA.xml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.path)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "parser:")
				assert.Contains(t, err.Error(), tt.path)
				if tt.name == "unsupported file type" {
					assert.True(t, errors.Is(err, ErrUnsupportedFileType))
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
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
