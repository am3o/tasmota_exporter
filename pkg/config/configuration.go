package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Devices []struct {
		Name     string `yaml:"Name"`
		IP       string `yaml:"IP"`
		Type     string `yaml:"Type"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"devices"`
}

func New(path string) (Configuration, error) {
	file, err := os.OpenFile(path, 0, 0)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not load configuration: %w", err)
	}

	var config Configuration
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Configuration{}, fmt.Errorf("could not unmarshal configuration file: %w", err)
	}

	return config, nil
}
