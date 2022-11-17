package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	JWT    string `json:"jwt"`
	Server string `json:"server"`
}

const DefaultPath = "$HOME/.config/dp-cli/config.json"

func Load(path string) (*Config, error) {
	configBytes, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save(path string) error {
	configBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, configBytes, 0600)
}
