package dag

import (
    "os"
    "fmt"
    "io"
    "encoding/json"
    "github.com/sirupsen/logrus"
)

// Use github.com/heimdalr/dag here
func LoadPipeline(filename string, l *logrus.Logger) (*Pipeline, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var pipeline Pipeline
	if err := json.Unmarshal(bytes, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

    pipeline.SetLogger(l)

	return &pipeline, nil
}

func LoadCatalog(filename string) (*map[string]PipelineFile, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    bytes, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    var catalog PipelineFileContainer
    if err := json.Unmarshal(bytes, &catalog); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    return &catalog.Catalog, nil
}
