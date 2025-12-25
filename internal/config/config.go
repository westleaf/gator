package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	bytes, _ := io.ReadAll(jsonFile)

	var config Config

	err = json.Unmarshal([]byte(bytes), &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(user string) error {
	conf, err := Read()
	if err != nil {
		return err
	}

	conf.CurrentUserName = user
	conf.DbURL = c.DbURL

	err = c.write(conf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, b, 0444)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s", homeDir, configFileName)
	return path, nil
}
