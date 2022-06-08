package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type HttpConfig struct {
	Port int `yaml:"port"`
}

type CacheConfig struct {
	ConnString string `yaml:"connection-string"`
}

type Config struct {
	HttpConfig  HttpConfig   `yaml:"http"`
	CacheConfig *CacheConfig `yaml:"cache,omitempty"`
}

func New(configPath string) (Config, error) {
	var data []byte
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	if os.Getenv("DEBUG") != "" {
		value, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Printf("%s", string(value))
		}

		fmt.Printf(string(value))
	}

	return config, nil
}
