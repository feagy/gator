package config

import (
    "os"
    "encoding/json"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DbUrl string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
    configPath, err := getConfigFilePath()
    if err != nil {
        return Config{}, err
    }

    configJson, err := os.Open(configPath)
    if err != nil {
		return Config{}, err
	}
	defer configJson.Close()

    config := Config{}
    jsonParser := json.NewDecoder(configJson)
    err = jsonParser.Decode(&config)
    if err != nil {
		return Config{}, err
	}
    return config, nil
}

func (c *Config) SetUser(user_name string) error {
    configPath, err := getConfigFilePath()
    if err != nil {
        return err
    }

    c.CurrentUserName = user_name
    
    configJson, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer configJson.Close()
    
	encoder := json.NewEncoder(configJson)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	return nil
}
