package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
				binaryName,
				filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileA.json"),
				filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileB.json"),
			},
		},
		{
			name:    "invalid args",
			args:    []string{binaryName, "fileA.json"},
			wantErr: "expected 2 file paths, got 1",
		},
		{
			name:    "gen diff error",
			args:    []string{binaryName, "nonexistent.json", "alsonothere.json"},
			wantErr: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args)
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

func TestMainError(t *testing.T) {
	if os.Getenv("TEST_GENDIFF_MAIN") == "1" {
		os.Args = []string{binaryName, "fileA.json"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainError")
	cmd.Env = append(os.Environ(), "TEST_GENDIFF_MAIN=1")

	err := cmd.Run()
	require.Error(t, err)

	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "expected process to exit with non-zero status")
	assert.Equal(t, 1, exitErr.ExitCode())
}
