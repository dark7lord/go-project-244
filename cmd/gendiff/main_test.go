package main

import (
	"os"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	oldArgs := os.Args

	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"gendiff",
		"../../testdata/fixture/fileA.json",
		"../../testdata/fixture/fileB.json",
	}

	err := Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
