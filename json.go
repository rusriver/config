package config

import (
	"encoding/json"
	"io/ioutil"
)

// ParseJson reads a JSON configuration from the given string.
func ParseJson(c string) (*Config, error) {
	return parseJson([]byte(c))
}

// ParseJsonFile reads a JSON configuration from the given filename.
func ParseJsonFile(filename string) (*Config, error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseJson(c)
}

// parseJson performs the real JSON parsing.
func parseJson(c []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = json.Unmarshal(c, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderJson renders a JSON configuration.
func RenderJson(c interface{}) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
