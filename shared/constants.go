package shared

import (
	"os"
	"path/filepath"
)

func GetConfigFolder() (string, error) {
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeFolder, ".config", "dag_doctor"), nil
}

func GetLogFilename() (string, error) {
	configFolder, err := GetConfigFolder()
	if err != nil {
		return "", nil
	}
	return filepath.Join(configFolder, "log.json"), nil
}

func GetOptFilename() (string, error) {
	configFolder, err := GetConfigFolder()
	if err != nil {
		return "", nil
	}
	return filepath.Join(configFolder, "opt.json"), nil
}

func GetDBFilename() (string, error) {
	configFolder, err := GetConfigFolder()
	if err != nil {
		return "", nil
	}
	return filepath.Join(configFolder, "data.db"), nil
}

func GetDAGFolder() (string, error) {
	configFolder, err := GetConfigFolder()
	if err != nil {
		return "", nil
	}
	return filepath.Join(configFolder, "dags"), nil
}
