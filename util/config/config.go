package config

import (
	"encoding/json"
	"errors"
	"os"
)

// errors introduced by the config package
var (
	ErrPathEmpty = errors.New("config: path is empty")
)

// ReadConfig reads the config file and unmarshals it into the given interface
func ReadConfig(configPath string, config interface{}) error {
	if configPath == "" {
		return ErrPathEmpty
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(configData, config)
}

// SaveConfig writes the given config to the given path
func SaveConfig(configPath string, config interface{}) error {
	if configPath == "" {
		return ErrPathEmpty
	}

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configData, 0644)
}
