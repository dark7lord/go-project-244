package code_test

import (
	"os"
	"strings"
	"testing"

	"code"
	"code/parser"
)


func TestGendiff(t *testing.T) {
	pathA := "testdata/file1.json"
	pathB := "testdata/file2.json"

	dataA, errA := parser.Parse(pathA)
	dataB, errB := parser.Parse(pathB)

	if errA != nil || errB != nil {
		t.Fatal("error parsing file")
	}
	
	got := strings.TrimSpace(code.GenDiff(dataA, dataB))

	expected, err := os.ReadFile("testdata/expected.txt")
	if err != nil {
		t.Fatal("error reading file with expected result")
	}
	want := strings.TrimSpace(string(expected))

	if want != got {
		t.Errorf("expected %s not equal actual %s", want, got)
	}
}