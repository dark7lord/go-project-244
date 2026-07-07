package formatters_test

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

func TestJSONFormatter(t *testing.T) {
	tests := []diffTestCase{
		{
			name: "nested node",
			diff: diff.Node{
				TypeDiff: diff.Nested,
				Children: []diff.Node{
					{
						Key:      "added",
						TypeDiff: diff.Added,
						NewValue: testNamePrimitiveValue,
					},
					{
						Key:      "nested",
						TypeDiff: diff.Nested,
						Children: []diff.Node{
							{
								Key:      "changed",
								TypeDiff: diff.Changed,
								OldValue: "old",
								NewValue: "new",
							},
						},
					},
				},
			},
			path: filepath.Join("testdata", "fixtures", "nestedNode.json"),
		},
		{
			name: testNameEmptyNested,
			diff: emptyNestedDiff,
			want: `{}`,
		},
	}
	runDiffTests(t, formatters.JSON, tests)
}

func TestJSONFormatterNoError(t *testing.T) {
	_, err := formatters.Format(string(formatters.JSON), diff.Node{
		TypeDiff: diff.Unchanged,
		OldValue: "test",
	})

	require.NoError(t, err)
}

func TestJSONFormatterUnknownKind(t *testing.T) {
	result, err := formatters.Format(string(formatters.JSON), diff.Node{TypeDiff: diff.Kind("bogus")})
	require.NoError(t, err)
	require.Equal(t, "{}", result)
}

func TestJSONFormatterUnknownChildKind(t *testing.T) {
	node := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{Key: "unknown", TypeDiff: diff.Kind("bogus")},
		},
	}

	expected := readFileToString(t, filepath.Join("testdata", "fixtures", "unknownChildKind.json"))
	result, err := formatters.Format(string(formatters.JSON), node)

	require.NoError(t, err)
	require.JSONEq(t, expected, result)
}

func TestJSONFormatterKeepsNilAndFalsyValues(t *testing.T) {
	node := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "addedNull",
				TypeDiff: diff.Added,
				NewValue: nil,
			},
			{
				Key:      "changedZeroToEmptyString",
				TypeDiff: diff.Changed,
				OldValue: 0,
				NewValue: "",
			},
			{
				Key:      "changedNullToFalse",
				TypeDiff: diff.Changed,
				OldValue: nil,
				NewValue: false,
			},
			{
				Key:      "removedNull",
				TypeDiff: diff.Removed,
				OldValue: nil,
			},
			{
				Key:      "unchangedNull",
				TypeDiff: diff.Unchanged,
				OldValue: nil,
			},
		},
	}

	expected := readFileToString(t, filepath.Join("testdata", "fixtures", "nilAndFalsyValues.json"))
	result, err := formatters.Format(string(formatters.JSON), node)

	require.NoError(t, err)
	require.JSONEq(t, expected, result)
}

func TestJSONFormatterMarshalError(t *testing.T) {
	node := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "nan",
				TypeDiff: diff.Added,
				NewValue: math.NaN(),
			},
		},
	}

	formatter, err := formatters.NewFormatter(string(formatters.JSON))
	require.NoError(t, err)

	_, err = formatter.Format(node)
	require.Error(t, err)
	require.Contains(t, err.Error(), "json marshal error")
}
