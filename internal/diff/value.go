package diff

// Value interface represents a value in JSON/YAML
type Value interface {
	isValue()
}

// Map represents a JSON/YAML object
type Map map[string]Value

// Slice represents a JSON/YAML array
type Slice []Value

// Number represents a JSON/YAML number
type Number float64

// String represents a JSON/YAML string
type String string

// Boolean represents a JSON/YAML boolean
type Boolean bool

// Null represents a JSON/YAML null
type Null struct{}

func (Map) isValue()     {}
func (Slice) isValue()   {}
func (Number) isValue()  {}
func (String) isValue()  {}
func (Boolean) isValue() {}
func (Null) isValue()    {}

// ToValue converts a native Go value (from JSON/YAML) to a Value
func ToValue(v any) Value {
	switch val := v.(type) {
	case map[string]any:
		m := make(Map, len(val))
		for k, v := range val {
			m[k] = ToValue(v)
		}

		return m
	case []any:
		s := make(Slice, len(val))
		for i, v := range val {
			s[i] = ToValue(v)
		}

		return s
	case float64:
		return Number(val)
	case string:
		return String(val)
	case bool:
		return Boolean(val)
	case nil:
		return Null{}
	case int:
		return Number(val)
	default:
		return Null{}
	}
}

// ToNative function converts a Value to a native Go type (map, slice, float64, string, bool, or nil)
func ToNative(v Value) any {
	switch val := v.(type) {
	case Map:
		m := make(map[string]any, len(val))
		for k, v := range val {
			m[k] = ToNative(v)
		}

		return m
	case Slice:
		s := make([]any, len(val))
		for i, v := range val {
			s[i] = ToNative(v)
		}

		return s
	case Number:
		return float64(val)
	case String:
		return string(val)
	case Boolean:
		return bool(val)
	case Null:
		return nil
	default:
		return nil
	}
}

// unknownValue satisfies the sealed Value interface for testing
// default branches in type switches.
type unknownValue struct{}

func (unknownValue) isValue() {}

// UnknownValue returns a Value that is not one of the standard types.
func UnknownValue() Value { return unknownValue{} }
