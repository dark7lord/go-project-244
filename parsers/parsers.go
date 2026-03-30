// Package parsers provides functions for parsing files (json, yml)
package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

// Parse function returns a json-like structure
func Parse(path string) (any, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if fileInfo.IsDir() {
		return "", fmt.Errorf("%q is directory", path)
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch ext := filepath.Ext(path); ext {
	case ".json":
		return parseJSON(fileBytes)
	case ".yml", ".yaml":
		return parseYAML(fileBytes)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}
