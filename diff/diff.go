// Package diff provided common type Diff for all packages
package diff

// Kind type represents the type of difference between two values
type Kind string

// Constants of differences between two values
const (
	Added     Kind = "added"
	Removed   Kind = "removed"
	Changed   Kind = "changed"
	Unchanged Kind = "unchanged"
)

// Diff struct represents the difference between two values
type Diff struct {
	Key      string
	TypeDiff Kind
	OldValue any
	NewValue any
}
