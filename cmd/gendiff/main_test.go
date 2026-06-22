package main

import (
	"os"
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

func TestMainExitsWithError(t *testing.T) {
	defer func() { osExit = os.Exit }()

	var got int
	osExit = func(code int) { got = code }

	os.Args = []string{binaryName, "fileA.json"}
	main()

	require.Equal(t, ExitError, got)
}

func TestMainExitsWithSuccess(t *testing.T) {
	defer func() { osExit = os.Exit }()

	var got int
	osExit = func(code int) { got = code }

	os.Args = []string{
		binaryName,
		filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileB.json"),
	}
	main()

	require.Equal(t, ExitOK, got)
}
