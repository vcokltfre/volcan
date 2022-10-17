package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config ConfigSchema

type ConfigSchema struct {
	Prefixes Prefixes `yaml:"prefixes"`
}

type Prefixes struct {
	Default string            `yaml:"default"`
	Servers map[string]string `yaml:"servers"`
}

func LoadConfig(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &Config)
}
