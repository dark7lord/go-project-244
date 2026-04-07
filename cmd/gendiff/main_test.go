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
		"gendiff",
		filepath.Join("..", "..", "testdata", "fixture", "fileA.json"),
		filepath.Join("..", "..", "testdata", "fixture", "fileB.json"),
	}

	err := Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
