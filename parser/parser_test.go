package parser_test

import (
	"code/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected any
		wantErr  bool
	}{
		{
			name:    "non existent file",
			path:    "non-existent.json",
			wantErr: true,
		},
		{
			name:    "directory",
			path:    ".",
			wantErr: true,
		},
		{
			name:    "invalid json",
			path:    "../testdata/fixture/invalid.json",
			wantErr: true,
		},
		{
			name: "valid json",
			path: "../testdata/fixture/fileA.json",
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.path)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
