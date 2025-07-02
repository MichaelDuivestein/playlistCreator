package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	InputPath          string   `json:"music folder"`
	OutputPath         string   `json:"output folder"`
	PlaylistName       string   `json:"playlist name"`
	ExtensionWhitelist []string `json:"Limit extensions to"`
	ShufflePlaylist    bool     `json:"shuffle playlist"`
	SplitPlaylist      bool     `json:"split playlist"`
	ChunkSize          int      `json:"split playlist into chunks of"`
	ListExtensions     bool     `json:"print a list of extensions"`
	ListFiles          bool     `json:"print a list of file names"`
	ListLimit          int      `json:"limit list file names to"`
	ReadTags           bool     `json:"read tags"`
}

func LoadConfig() (*Config, error) {
	return loadConfigFromFile("config/config.json")
}

func loadConfigFromFile(filePathAndName string) (*Config, error) {
	data, err := os.ReadFile(filePathAndName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.InputPath == "" {
		return nil, errors.New("config.InputPath is empty")
	}
	if config.OutputPath == "" {
		return nil, errors.New("config.OutputPath is empty")
	}

	return &config, nil
}
