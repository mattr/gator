package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	CurrentUserName string `json:"current_user_name"`
	DatabaseURL     string `json:"db_url"`
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(cfg)
}

func Read() (Config, error) {
	var config Config
	filepath, err := getConfigFilePath()
	if err != nil {
		return config, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFileName, nil
}

func write(cfg *Config) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
