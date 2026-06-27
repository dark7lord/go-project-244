// Package parsers provides functions for parsing files (json, yml)
package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ErrUnsupportedFileType is returned when the file type is not supported
var ErrUnsupportedFileType = errors.New("unsupported file type")

// ErrUnsupportedRootType is returned when the parsed file root is not an object.
var ErrUnsupportedRootType = errors.New("unsupported root type: expected map[string]any")

func parseJSON(data []byte) (map[string]any, error) {
	var result map[string]any

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func parseYAML(data []byte) (map[string]any, error) {
	var result map[string]any

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	normalized := normalizeYAML(result)
	m, ok := normalized.(map[string]any)
	if !ok {
		return nil, ErrUnsupportedRootType
	}

	return m, nil
}

func normalizeYAML(v any) any {
	switch val := v.(type) {
	case int:
		return float64(val)
	case map[string]any:
		m := make(map[string]any, len(val))
		for k, v := range val {
			m[k] = normalizeYAML(v)
		}

		return m
	case []any:
		s := make([]any, len(val))
		for i, v := range val {
			s[i] = normalizeYAML(v)
		}

		return s
	default:
		return v
	}
}

func ensureMap(v any) (map[string]any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, ErrUnsupportedRootType
	}

	return m, nil
}

func parseError(path, ext, op string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("parse error: path=%q type=%q %s: %w", path, ext, op, err)
}

// Parse function returns a parsed JSON/YAML structure
func Parse(path string) (map[string]any, error) {
	ext := filepath.Ext(path)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, parseError(path, ext, "stat", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("parse error: path=%q is directory", path)
	}

	if ext == "" {
		return nil, fmt.Errorf("parse error: in path=%q missing file extension", path)
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, parseError(path, ext, "read", err)
	}

	switch ext {
	case ".json":
		value, err := parseJSON(fileBytes)
		if err != nil {
			return nil, parseError(path, ext, "parse", err)
		}

		return ensureMap(value)
	case ".yml", ".yaml":
		value, err := parseYAML(fileBytes)
		if err != nil {
			return nil, parseError(path, ext, "parse", err)
		}

		return ensureMap(value)
	default:
		return nil, fmt.Errorf("parser: path=%q type=%q: %w", path, ext, ErrUnsupportedFileType)
	}
}
