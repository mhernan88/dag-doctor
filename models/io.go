package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadDAG(filename string) (*DAG, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	dagContainer := make(map[string]map[string]Node)
	if err := json.Unmarshal(bytes, &dagContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	nodes, ok := dagContainer["nodes"]
	if !ok {
		return nil, fmt.Errorf("'nodes' key not found in data")
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("0 nodes found in catalog")
	}

	dag, err := NewDAGFromMap(nodes)
	if err != nil {
		return nil, err
	}

	return dag, nil
}

func SaveState(filename string, state *State) error {
	f, err := os.OpenFile(filename, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(state)
}

func LoadState(filename string) (*State, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var state State
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}
