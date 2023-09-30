package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/models"
)

func PrintTree(title string, content string) {
	fmt.Printf("%s\n", title)
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
	os.MkdirAll(folder, 0777)

	fileContents, err := os.ReadFile(dagFilename)
	if err != nil {
		return "", fmt.Errorf("failed to read dag | %v", err)
	}
	
	outputFilename := filepath.Join(folder, fmt.Sprintf("%s.json", dagID))
	err = os.WriteFile(outputFilename, fileContents, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to write dag | %v", err)
	}
	return outputFilename, nil
}

func SaveStateToRepo(state *models.State) (string, error) {
	folder, err := GetStateFolder()
	if err != nil {
		return "", fmt.Errorf("failed to acquire state folder | %v", err)
	}
	os.MkdirAll(folder, 0777)

	filename := filepath.Join(folder, fmt.Sprintf("%s.json", uuid.NewString()))
	models.SaveState(filename, state)
	return filename, nil
}
