package formatters_test

import "code/internal/diff"

const (
	testNameNestedAll      = "nested all types"
	testNameEmptyNested    = "empty nested object"
	testNamePrimitiveValue = "value"
)

const singleAddedNodeJSON = `{
		"type": "nested",
		"children": [
			{
				"key": "Key",
				"type": "added",
				"newValue": "value"
			}
		]
	}`

func singleAddedDiff() diff.Node {
	return diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "Key",
				TypeDiff: diff.Added,
				NewValue: testNamePrimitiveValue,
			},
		},
	}
}

var flatDiffAll = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{
		{Key: "stringKey", TypeDiff: diff.Unchanged, OldValue: "string value"},
		{Key: "numberKey", TypeDiff: diff.Unchanged, OldValue: 42.0},
		{Key: "boolKey", TypeDiff: diff.Unchanged, OldValue: true},
		{Key: "nullKey", TypeDiff: diff.Unchanged, OldValue: nil},
		{Key: "addedKey", TypeDiff: diff.Added, NewValue: "added value"},
		{Key: "removedKey", TypeDiff: diff.Removed, OldValue: "removed value"},
		{Key: "changedKey", TypeDiff: diff.Changed, OldValue: "old value", NewValue: "new value"},
	},
}

var nestedDiffAll = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{
		{
			Key:      "outerUnchanged",
			TypeDiff: diff.Unchanged,
			OldValue: "value1",
		},
		{
			Key:      "outerNested",
			TypeDiff: diff.Nested,
			Children: []diff.Node{
				{
					Key:      "nestedUnchanged",
					TypeDiff: diff.Unchanged,
					OldValue: "nestedValue",
				},
				{
					Key:      "nestedAdded",
					TypeDiff: diff.Added,
					NewValue: "addedNestedValue",
				},
				{
					Key:      "nestedRemoved",
					TypeDiff: diff.Removed,
					OldValue: "removedNestedValue",
				},
				{
					Key:      "nestedDeep",
					TypeDiff: diff.Nested,
					Children: []diff.Node{
						{
							Key:      "deepChanged",
							TypeDiff: diff.Changed,
							OldValue: "oldDeepValue",
							NewValue: "newDeepValue",
						},
					},
				},
			},
		},
		{
			Key:      "outerAdded",
			TypeDiff: diff.Added,
			NewValue: map[string]any{"addedNestedKey": "addedNestedValue"},
		},
		{
			Key:      "outerRemoved",
			TypeDiff: diff.Removed,
			OldValue: map[string]any{"removedNestedKey": "removedNestedValue"},
		},
		{
			Key:      "outerChanged",
			TypeDiff: diff.Changed,
			OldValue: map[string]any{"oldNestedKey": "oldNestedValue"},
			NewValue: map[string]any{"newNestedKey": "newNestedValue"},
		},
	},
}

var emptyNestedDiff = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{},
}
