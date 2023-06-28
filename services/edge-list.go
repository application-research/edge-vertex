package services

import (
	"encoding/json"
	"os"
)

func ReadEdgeList(filename string) ([]string, error) {
	// Read the file from disk
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Unmarshal the file contents into a slice of strings
	var edges []string
	err = json.NewDecoder(file).Decode(&edges)
	if err != nil {
		return nil, err
	}

	return edges, nil
}
