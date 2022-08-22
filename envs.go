package config

import (
	"strings"
	"syscall"
)

// Fetch data from system env, based on existing config keys.
func (c *Config) Env() *Config {
	return c.EnvPrefix("")
}

// Fetch data from system env using prefix, based on existing config keys.
func (c *Config) EnvPrefix(prefix string) *Config {
	if prefix != "" {
		prefix = strings.ToUpper(prefix) + "_"
	}

	keys := getKeys(c.Root)
	for _, key := range keys {
		k := strings.ToUpper(strings.Join(key, "_"))
		if val, exist := syscall.Getenv(prefix + k); exist {
			c.Set(strings.Join(key, "."), val)
		}
	}
	return c
}
