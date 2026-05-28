package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSuccess(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = writePipe

	os.Args = []string{
		binaryName,
		filepath.Join("..", "..", "testdata", "fixture", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixture", "fileB.json"),
	}

	err = Run()
	require.NoError(t, writePipe.Close())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("expected diff output, got empty string")
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

func TestRunUnsupportedFormat(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		binaryName,
		"--format",
		"invalid",
		filepath.Join("..", "..", "testdata", "fixture", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixture", "fileB.json"),
	}

	err := Run()
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}

	expected := "unsupported format: invalid"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestMainCallsRun(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = writePipe

	os.Args = []string{
		binaryName,
		filepath.Join("..", "..", "testdata", "fixture", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixture", "fileB.json"),
	}

	main()
	require.NoError(t, writePipe.Close())

	output, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("expected diff output from main, got empty string")
	}
}
