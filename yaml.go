package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ParseYamlBytes reads a YAML configuration from the given []byte.
func ParseYamlBytes(c []byte) (*Config, error) {
	return parseYaml(c)
}

// ParseYaml reads a YAML configuration from the given string.
func ParseYaml(c string) (*Config, error) {
	return parseYaml([]byte(c))
}

// ParseYamlFile reads a YAML configuration from the given filename.
func ParseYamlFile(filename string) (*Config, error) {
	c, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseYaml(c)
}

// parseYaml performs the real YAML parsing.
func parseYaml(c []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = yaml.Unmarshal(c, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderYaml renders a YAML configuration.
func RenderYaml(c interface{}) (string, error) {
	b, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
