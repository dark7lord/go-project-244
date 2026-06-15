package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMapKey = "key"

func TestToNative(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  any
	}{
		{
			name:  "map",
			input: Map{testMapKey: String("val")},
			want:  map[string]any{testMapKey: "val"},
		},
		{
			name:  "slice",
			input: Slice{Number(1), Number(2)},
			want:  []any{float64(1), float64(2)},
		},
		{
			name:  "number",
			input: Number(42.5),
			want:  float64(42.5),
		},
		{
			name:  "string",
			input: String("hello"),
			want:  "hello",
		},
		{
			name:  "bool",
			input: Boolean(true),
			want:  true,
		},
		{
			name:  "null",
			input: Null{},
			want:  nil,
		},
		{
			name:  "nested map",
			input: Map{"a": Map{"b": Number(1)}},
			want:  map[string]any{"a": map[string]any{"b": float64(1)}},
		},
		{
			name:  "empty slice",
			input: Slice{},
			want:  []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToNative(tt.input))
		})
	}
}

func TestToValue(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  Value
	}{
		{
			name:  "map",
			input: map[string]any{testMapKey: "val"},
			want:  Map{testMapKey: String("val")},
		},
		{
			name:  "slice",
			input: []any{1.0, 2.0},
			want:  Slice{Number(1), Number(2)},
		},
		{
			name:  "float64",
			input: float64(42.5),
			want:  Number(42.5),
		},
		{
			name:  "string",
			input: "hello",
			want:  String("hello"),
		},
		{
			name:  "bool",
			input: true,
			want:  Boolean(true),
		},
		{
			name:  "nil",
			input: nil,
			want:  Null{},
		},
		{
			name:  "int",
			input: 42,
			want:  Number(42),
		},
		{
			name:  "nested map",
			input: map[string]any{"a": map[string]any{"b": 1.0}},
			want:  Map{"a": Map{"b": Number(1)}},
		},
		{
			name:  "empty slice",
			input: []any{},
			want:  Slice{},
		},
		{
			name:  "unknown type",
			input: struct{}{},
			want:  Null{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToValue(tt.input))
		})
	}
}

type testUnknownValue struct{}

func (testUnknownValue) isValue() {}

func TestToNativeDefault(t *testing.T) {
	assert.Equal(t, nil, ToNative(testUnknownValue{}))
}
