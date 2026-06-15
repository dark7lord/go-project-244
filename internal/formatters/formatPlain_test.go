package formatters_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestFormatPlainNoTrailingNewline(t *testing.T) {
	diffNode := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "addedKey",
				TypeDiff: diff.Added,
				NewValue: diff.String("added value"),
			},
		},
	}

	got, err := formatters.PrintDiff(diffNode, formatters.Plain)
	require.NoError(t, err, "PrintDiff() returned an error: %v", err)
	require.Equal(t, "Property 'addedKey' was added with value: 'added value'", got)
}

func TestFormatPlain(t *testing.T) {
	tests := []diffTestCase{
		{
			name: "flat all types",
			diff: flatDiffAll,
			path: filepath.Join("testdata", "fixtures", "plainFlat.txt"),
		},
		{
			name: testNameNestedAll,
			diff: nestedDiffAll,
			path: filepath.Join("testdata", "fixtures", "plainNested.txt"),
		},
		{
			name: testNameEmptyNested,
			diff: emptyNestedDiff,
			want: "",
		},
		{
			name: "added complex value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{Key: "complexKey", TypeDiff: diff.Added, NewValue: diff.Slice{diff.Number(1), diff.Number(2)}},
				},
			},
			want: "Property 'complexKey' was added with value: [complex value]",
		},
		{
			name: "added nil value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{Key: "nullKey", TypeDiff: diff.Added, NewValue: diff.Null{}},
				},
			},
			want: "Property 'nullKey' was added with value: null",
		},
		{
			name: "added bool value",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{Key: "boolKey", TypeDiff: diff.Added, NewValue: diff.Boolean(true)},
				},
			},
			want: "Property 'boolKey' was added with value: true",
		},
	}
	runDiffTests(t, formatters.Plain, tests)
}
