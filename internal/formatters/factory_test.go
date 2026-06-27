package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/formatters"
)

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: string(formatters.Stylish), want: true},
		{name: string(formatters.Plain), want: true},
		{name: string(formatters.JSON), want: true},
		{name: "xml", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, formatters.IsValidFormat(tt.name))
		})
	}
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format formatters.OutputFormat
		want   string
	}{
		{
			name:   string(formatters.Stylish),
			format: formatters.Stylish,
			want:   "{\n  + Key: value\n}",
		},
		{
			name:   string(formatters.Plain),
			format: formatters.Plain,
			want:   "Property 'Key' was added with value: 'value'",
		},
		{
			name:   string(formatters.JSON),
			format: formatters.JSON,
			want:   singleAddedNodeJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := formatters.NewFormatter(string(tt.format))
			require.NoError(t, err)

			got, err := formatter.Format(singleAddedDiff())
			require.NoError(t, err)
			requireOutputEqual(t, tt.format, tt.want, got)
		})
	}
}

func TestNewFormatterUnsupportedFormat(t *testing.T) {
	_, err := formatters.NewFormatter("xml")
	require.Error(t, err)
}
