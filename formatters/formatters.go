// Package formatters provides functions for formatting the difference between files
package formatters

// PrintDiff function returns the difference between two structures as a string in the specified format
func PrintDiff(diff any, format string) string {
	switch format {
	case "stylish":
		return printStylishDiff(diff, 0)
	case "json":
		return printJSONDiff(diff)
	default:
		return printPlainDiff(diff, []string{})
	}
}
