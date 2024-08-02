package config

import (
	"encoding/json"
	"os"
	"time"
)

const perm = 0644

type Config struct {
	Name                  string
	Filepath              string
	TelegramApiToken      string
	Database              string
	DelayBetweenSnapshots time.Duration
}

func LoadFromFile(filepath string) (*Config, error) {
	c := &Config{Filepath: filepath}

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			data, err = json.Marshal(c)
			if err != nil {
				return nil, err
			}
			os.WriteFile(filepath, data, perm)
			return c, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	c.Save()

	return c, nil
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.Filepath, data, perm)
	if err != nil {
		return err
	}

	return nil
}
