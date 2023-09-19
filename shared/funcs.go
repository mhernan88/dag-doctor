package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PrintTree(title string, content string) {
	fmt.Printf("%s:\n", title)
	fmt.Println(strings.Repeat("-", len(title)+1))
	for _, line := range strings.Split(content, "\n") {
		fmt.Printf("%s\n", line)
	}
	fmt.Println()
}

func CopyDAGToRepo(dagFilename string, dagID string) (string, error) {
	folder, err := GetDAGFolder()
	if err != nil {
		return "", fmt.Errorf("failed to acquire dag folder | %v", err)
	}
	os.MkdirAll(folder, 0666)

	fileContents, err := os.ReadFile(dagFilename)
	if err != nil {
		return "", fmt.Errorf("failed to read dag | %v", err)
	}
	
	outputFilename := filepath.Join(folder, fmt.Sprintf("%s.json", dagID))
	return outputFilename, os.WriteFile(outputFilename, fileContents, 0666)
}
