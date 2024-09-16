package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PathConfig struct {
	LogPath string `json:"log_path"`
}

func getConfigPath() string {
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}
	return filepath.Join(filepath.Dir(execPath), "mb2wclfixer_path.json")
}

func SavePath(path string) error {
	config := PathConfig{LogPath: path}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

func LoadPath() (string, error) {
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var config PathConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return "", err
	}
	return config.LogPath, nil
}
