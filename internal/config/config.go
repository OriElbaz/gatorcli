package config

import (
	"fmt"
	"os"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Error with getting config file path: %v", err)
		return Config{}, err
	}

	jsonData, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error with JSONDATA: %v", err)
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		fmt.Printf("Error with UNMARSHAL: %v\n", err)
		return Config{}, err
	}

	return config, nil
}

func (c Config) SetUser(username string) error {
	c.CurrentUserName = username

	err := write(c)
	if err != nil {
		fmt.Printf("ERROR writing to config file: %v\n", err)
	}

	return nil
}

func write(cfg Config) error {
	jsonData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		fmt.Printf("ERROR converting config struct to JSON: %v\n", err)
	}

	outputPath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("ERROR getting config file path: %v\n", err)
		return err
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("ERROR writing to config file: %v\n", err)
		return err
	}

	return nil

}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error with getting config file path: %v\n", err)
		return "", err
	}

	return homePath + "/" + configFileName, nil
}

/*
Learnt:
- os.HomeDir
- os.ReadFile
- os.WriteFile
- json.MarshalIndent -> struct to JSON

*/