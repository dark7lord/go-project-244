package formatters_test

import (
	"math"
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
			want: `{
				"type": "nested",
				"children": [
					{
						"key": "added",
						"type": "added",
						"newValue": "value"
					},
					{
						"key": "nested",
						"type": "nested",
						"children": [
							{
								"key": "changed",
								"type": "changed",
								"oldValue": "old",
								"newValue": "new"
							}
						]
					}
				]
			}`,
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

	result, err := formatters.Format(string(formatters.JSON), node)

	require.NoError(t, err)
	require.JSONEq(t, `{
		"type": "nested",
		"children": [
			{
				"key": "unknown",
				"type": "bogus"
			}
		]
	}`, result)
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

	result, err := formatters.Format(string(formatters.JSON), node)

	require.NoError(t, err)
	require.JSONEq(t, `{
		"type": "nested",
		"children": [
			{
				"key": "addedNull",
				"type": "added",
				"newValue": null
			},
			{
				"key": "changedZeroToEmptyString",
				"type": "changed",
				"oldValue": 0,
				"newValue": ""
			},
			{
				"key": "changedNullToFalse",
				"type": "changed",
				"oldValue": null,
				"newValue": false
			},
			{
				"key": "removedNull",
				"type": "removed",
				"oldValue": null
			},
			{
				"key": "unchangedNull",
				"type": "unchanged",
				"oldValue": null
			}
		]
	}`, result)
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
