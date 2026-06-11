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
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		binaryName,
		filepath.Join("..", "..", "testdata", "gendiff", "fixture", "flat", "fileA.json"),
		filepath.Join("..", "..", "testdata", "gendiff", "fixture", "flat", "fileB.json"),
	}

	err := Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunInvalidArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{binaryName, "fileA.json"}

	err := Run()
	if err == nil {
		t.Fatal("expected error for invalid args count, got nil")
	}

	expected := "expected 2 file paths, got 1"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestRunGenDiffError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{binaryName, "nonexistent.json", "alsonothere.json"}

	err := Run()
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
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
