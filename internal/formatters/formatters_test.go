package formatters_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters"
)

const (
	stylish = "stylish"
	plain   = "plain"
	json    = "json"
)

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   bool
	}{
		{name: stylish, format: stylish, want: true},
		{name: plain, format: plain, want: true},
		{name: json, format: json, want: true},
		{name: "unknown", format: "xml", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, formatters.IsValidFormat(tt.format))
		})
	}
}

func TestPrintDiff(t *testing.T) {
	testDiff := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "Key",
				TypeDiff: diff.Added,
				NewValue: diff.String("value"),
			},
		},
	}

	tests := []struct {
		name   string
		format formatters.PrintFormat
		want   string
	}{
		{
			name:   stylish,
			format: formatters.Stylish,
			want:   "{\n  + Key: value\n}",
		},
		{
			name:   plain,
			format: formatters.Plain,
			want:   "Property 'Key' was added with value: 'value'",
		},
		{
			name:   json,
			format: formatters.JSON,
			want:   "[\n  {\n    \"key\": \"Key\",\n    \"type\": \"added\",\n    \"value\": \"value\"\n  }\n]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(testDiff, tt.format)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormat(t *testing.T) {
	testDiff := diff.Node{
		TypeDiff: diff.Nested,
		Children: []diff.Node{
			{
				Key:      "Key",
				TypeDiff: diff.Added,
				NewValue: diff.String("value"),
			},
		},
	}

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{name: stylish, format: stylish, want: "{\n  + Key: value\n}"},
		{name: plain, format: plain, want: "Property 'Key' was added with value: 'value'"},
		{name: json, format: json, want: "[\n  {\n    \"key\": \"Key\",\n    \"type\": \"added\",\n    \"value\": \"value\"\n  }\n]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.Format(tt.format, testDiff)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormatUnsupported(t *testing.T) {
	_, err := formatters.Format("xml", diff.Node{})
	require.Error(t, err)
}

func TestPrintDiffUnsupportedFormat(t *testing.T) {
	_, err := formatters.PrintDiff(diff.Node{}, formatters.PrintFormat("xml"))
	require.Error(t, err)
}

func TestFormatJSONMarshalError(t *testing.T) {
	node := diff.Node{
		TypeDiff: diff.Added,
		NewValue: diff.Number(math.NaN()),
	}

	_, err := formatters.PrintDiff(node, formatters.JSON)
	require.Error(t, err)
	require.Contains(t, err.Error(), "json marshal error")
}
