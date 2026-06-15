package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSuccess(t *testing.T) {
	err := run([]string{
		binaryName,
		filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixtures", "flat", "fileB.json"),
	})
	require.NoError(t, err)
}

func TestRunInvalidArgs(t *testing.T) {
	err := run([]string{binaryName, "fileA.json"})
	require.EqualError(t, err, "expected 2 file paths, got 1")
}

func TestRunGenDiffError(t *testing.T) {
	err := run([]string{binaryName, "nonexistent.json", "alsonothere.json"})
	require.Error(t, err)
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
