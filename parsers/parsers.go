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

var ErrUnsupportedFileType = errors.New("unsupported file type")

func parseJSON(data []byte) (any, error) {
	var result any

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func parseYAML(data []byte) (any, error) {
	var result any

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func parseError(path, ext, op string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("parser: path=%q type=%q %s: %w", path, ext, op, err)
}

// Parse function returns a json-like structure
func Parse(path string) (any, error) {
	ext := filepath.Ext(path)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, parseError(path, ext, "stat", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("parser: path=%q type=%q is directory", path, ext)
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, parseError(path, ext, "read", err)
	}

	switch ext {
	case ".json":
		value, err := parseJSON(fileBytes)
		return value, parseError(path, ext, "parse", err)
	case ".yml", ".yaml":
		value, err := parseYAML(fileBytes)
		return value, parseError(path, ext, "parse", err)
	default:
		return nil, parseError(path, ext, "unsupported file type", ErrUnsupportedFileType)
	}
}
