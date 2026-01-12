package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read turns config json file into config struct
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("get config file path: %w", err)
	}

	jsonDataBytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file to bytes: %w", err)
	}

	var config Config

	if err = json.Unmarshal(jsonDataBytes, &config); err != nil {
		return Config{}, fmt.Errorf("unmarshal to config struct: %w", err)
	}

	return config, nil
}

// SetUser updates config file username
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	err := write(c)
	if err != nil {
		return fmt.Errorf("write to config: %w", err)
	}

	return nil
}

func write(c *Config) error {
	jsonDataBytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("marshal config to bytes: %w", err)
	}

	outputPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("get config file path: %w", err)
	}

	if err = os.WriteFile(outputPath, jsonDataBytes, 0644); err != nil {
		return fmt.Errorf("write to config file: %w", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get config file path: %w", err)
	}

	return filepath.Join(homePath, configFileName), nil
}