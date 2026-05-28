package formatters_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"code/diff"
	"code/formatters"
)

func readFileToString(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file %s", path)

	return strings.TrimSpace(string(data))
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   bool
	}{
		{name: "stylish", format: "stylish", want: true},
		{name: "plain", format: "plain", want: true},
		{name: "json", format: "json", want: true},
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
				Key:      "key",
				TypeDiff: diff.Added,
				NewValue: "value",
			},
		},
	}

	tests := []struct {
		name    string
		format  formatters.PrintFormat
		want    string
		wantErr bool
	}{
		{
			name:   "stylish",
			format: formatters.Stylish,
			want:   "{\n  + key: value\n}",
		},
		{
			name:   "plain",
			format: formatters.Plain,
			want:   "Property 'key' was added with value: 'value'",
		},
		{
			name:   "json",
			format: formatters.JSON,
			want:   "{\n  \"key [added]\": \"value\"\n}",
		},
		{
			name:    "unsupported format",
			format:  formatters.PrintFormat("xml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatters.PrintDiff(testDiff, tt.format)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
