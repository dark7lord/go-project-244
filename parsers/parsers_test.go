package parsers_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/parsers"
)

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
			path:    "../testdata/fixture/invalid.json",
			wantErr: true,
		},
		{
			name:    "invalid yml",
			path:    "../testdata/fixture/invalid.yaml",
			wantErr: true,
		},
		{
			name: "valid json",
			path: "../testdata/fixture/fileA.json",
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name: "valid yml",
			path: "../testdata/fixture/fileA.yml",
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name:    "unsupported file type",
			path:    "../testdata/fixture/fileA.xml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsers.Parse(tt.path)

			if tt.wantErr {
				require.Error(t, err)
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

	_, err = parsers.Parse(path)
	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrPermission))
}
