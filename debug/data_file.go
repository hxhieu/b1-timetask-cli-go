package debug

import (
	"encoding/json"
	"os"
)

func WriteDataFile[T any](file string, data T) {
	if debugData, err := json.MarshalIndent(data, "", "\t"); err == nil {
		os.WriteFile(file, debugData, 0644)
	}
}

func LoadDataFile[T any](file string) *T {
	if buf, err := os.ReadFile(file); err == nil {
		result := new(T)
		if err = json.Unmarshal(buf, result); err == nil {
			return result
		}
	}
	return nil
}
