package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config ConfigSchema

type ConfigSchema struct {
	Prefixes      Prefixes       `yaml:"prefixes"`
	Levels        map[string]int `yaml:"levels"`
	CommandLevels map[string]int `yaml:"command_levels"`
	Modules       Modules        `yaml:"modules"`
}

func (c *ConfigSchema) IsEnabled(module string) bool {
	switch module {
	case "meta":
		return c.Modules.MetaModule.Enabled
	default:
		return false
	}
}

type Prefixes struct {
	Default string            `yaml:"default"`
	Servers map[string]string `yaml:"servers"`
}

type Modules struct {
	MetaModule MetaModule `yaml:"meta"`
}

type MetaModule struct {
	Enabled bool `yaml:"enabled"`
}

func LoadConfig(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &Config)
}
