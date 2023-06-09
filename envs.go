package config

import (
	"strings"
	"syscall"
)

// Fetch data from system env using prefix, based on existing config keys.
func (c *Config) ExtendByEnvs_WithPrefix(prefix string) *Config {
	if prefix != "" {
		prefix = strings.ToUpper(prefix) + "_"
	}
	keys := getAllKeys(c.Root)
	for _, key := range keys {
		k := strings.ReplaceAll(strings.ToUpper(strings.Join(key, "_")), "-", "")
		if val, exist := syscall.Getenv(prefix + k); exist {
			c.Set(strings.Join(key, "."), val)
		}
	}
	return c
}
