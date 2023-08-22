package data

import (
    "io"
    "os"
    "fmt"
    "encoding/json"
)

func LoadDAG(filename string)(map[string]*Node, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    bytes, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    dag := make(map[string]map[string]*Node)
    if err := json.Unmarshal(bytes, &dag); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    nodes, ok := dag["nodes"]
    if !ok {
        return nil, fmt.Errorf("'nodes' key not found in catalog")
    }

    return nodes, nil
}

func LoadCatalog(filename string) (map[string]Dataset, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    bytes, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    catalog := make(map[string]map[string]Dataset)
    if err := json.Unmarshal(bytes, &catalog); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    datasets, ok := catalog["datasets"]
    if !ok {
        return nil, fmt.Errorf("'datasets' key not found in catalog")
    }

    return datasets, nil
}
