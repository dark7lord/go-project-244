package formatters_test

import "code/internal/diff"

const (
	testNameNestedAll      = "nested all types"
	testNameEmptyNested    = "empty nested object"
	testNamePrimitiveValue = "value"
)

var flatDiffAll = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{
		{Key: "stringKey", TypeDiff: diff.Unchanged, OldValue: diff.String("string value")},
		{Key: "numberKey", TypeDiff: diff.Unchanged, OldValue: diff.Number(42)},
		{Key: "boolKey", TypeDiff: diff.Unchanged, OldValue: diff.Boolean(true)},
		{Key: "nullKey", TypeDiff: diff.Unchanged, OldValue: diff.Null{}},
		{Key: "addedKey", TypeDiff: diff.Added, NewValue: diff.String("added value")},
		{Key: "removedKey", TypeDiff: diff.Removed, OldValue: diff.String("removed value")},
		{Key: "changedKey", TypeDiff: diff.Changed, OldValue: diff.String("old value"), NewValue: diff.String("new value")},
	},
}

var nestedDiffAll = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{
		{
			Key:      "outerUnchanged",
			TypeDiff: diff.Unchanged,
			OldValue: diff.String("value1"),
		},
		{
			Key:      "outerNested",
			TypeDiff: diff.Nested,
			Children: []diff.Node{
				{
					Key:      "nestedUnchanged",
					TypeDiff: diff.Unchanged,
					OldValue: diff.String("nestedValue"),
				},
				{
					Key:      "nestedAdded",
					TypeDiff: diff.Added,
					NewValue: diff.String("addedNestedValue"),
				},
				{
					Key:      "nestedRemoved",
					TypeDiff: diff.Removed,
					OldValue: diff.String("removedNestedValue"),
				},
				{
					Key:      "nestedDeep",
					TypeDiff: diff.Nested,
					Children: []diff.Node{
						{
							Key:      "deepChanged",
							TypeDiff: diff.Changed,
							OldValue: diff.String("oldDeepValue"),
							NewValue: diff.String("newDeepValue"),
						},
					},
				},
			},
		},
		{
			Key:      "outerAdded",
			TypeDiff: diff.Added,
			NewValue: diff.Map{"addedNestedKey": diff.String("addedNestedValue")},
		},
		{
			Key:      "outerRemoved",
			TypeDiff: diff.Removed,
			OldValue: diff.Map{"removedNestedKey": diff.String("removedNestedValue")},
		},
		{
			Key:      "outerChanged",
			TypeDiff: diff.Changed,
			OldValue: diff.Map{"oldNestedKey": diff.String("oldNestedValue")},
			NewValue: diff.Map{"newNestedKey": diff.String("newNestedValue")},
		},
	},
}

var emptyNestedDiff = diff.Node{
	TypeDiff: diff.Nested,
	Children: []diff.Node{},
}
