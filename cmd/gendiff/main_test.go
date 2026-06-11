package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		binaryName,
		filepath.Join("..", "..", "testdata", "fixture", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixture", "fileB.json"),
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
	oldArgs := os.Args
	oldOsExit := osExit

	defer func() {
		os.Args = oldArgs
		osExit = oldOsExit
	}()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	os.Args = []string{binaryName, "fileA.json"}
	main()

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}
