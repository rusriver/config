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
	paths := getAllPaths(c.DataTreeRoot)
	for _, pathParts := range paths {
		k := strings.ReplaceAll(strings.ToUpper(strings.Join(pathParts, "_")), "-", "")
		if val, exist := syscall.Getenv(prefix + k); exist {
			c.Set(pathParts, val)
		}
	}
	return c
}
