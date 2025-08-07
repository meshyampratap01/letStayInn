package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJson(filepath string, v interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		data = []byte("[]")
	}
	return json.Unmarshal(data, v)
}

func WriteJson(filepath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return os.WriteFile(filepath, data, 0644)
}
