// Package diff provided common type Diff for all packages
package diff

// Diff struct represents the difference between two values
type Diff struct {
	TypeDiff string
	OldValue any
	NewValue any
}

// DiffTypes constants for all states Diff
const (
	DiffTypeAdded     = "added"
	DiffTypeRemoved   = "removed"
	DiffTypeChanged   = "changed"
	DiffTypeUnchanged = "unchanged"
)
