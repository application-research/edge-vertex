package services

import (
	"encoding/json"
	"os"
)

type ListReader interface {
	Read() ([]string, error)
}

type JsonEdgeListReader struct{ Filename string }

func (r JsonEdgeListReader) Read() ([]string, error) {
	// Read the file from disk
	file, err := os.Open(r.Filename)
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
