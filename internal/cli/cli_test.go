package cli

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name: "success",
			args: []string{
				BinaryName,
				filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileA.json"),
				filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileB.json"),
			},
		},
		{
			name:    "invalid args",
			args:    []string{BinaryName, "fileA.json"},
			wantErr: "expected 2 file paths, got 1",
		},
		{
			name:    "gen diff error",
			args:    []string{BinaryName, "nonexistent.json", "alsonothere.json"},
			wantErr: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args)
			switch tt.wantErr {
			case "":
				require.NoError(t, err)
			case "*":
				require.Error(t, err)
			default:
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestExecuteExitsWithError(t *testing.T) {
	var stderr bytes.Buffer

	got := Execute([]string{"gendiff", "fileA.json"}, &stderr)

	require.Equal(t, ExitError, got)
	require.Contains(t, stderr.String(), "expected 2 file paths, got 1")
}
