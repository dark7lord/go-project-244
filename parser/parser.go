package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// - некие заметки насчет путей (абсолютные и относительные пути)
// - определить формат по расширению файла

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

	var data any

	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}
