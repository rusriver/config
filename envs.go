package config

import (
	"os"
	"strings"
	"syscall"
)

// Fetch data from system env using prefix, based on existing config keys.
// VERY IMPORTANT USAGE NOTE: this can override what is already present in the config,
// but it cannot create new things, which were not in the config.
func (c *Config) ExtendByEnvs_WithPrefix(prefix string) *Config {
	if prefix != "" {
		prefix = strings.ToUpper(prefix) + "_"
	}
	paths := getAllPaths(c.DataSubTree)
	for _, pathParts := range paths {
		k := strings.ReplaceAll(strings.ToUpper(strings.Join(pathParts, "_")), "-", "")
		if val, exist := syscall.Getenv(prefix + k); exist {
			c.Set(pathParts, val)
		}
	}
	return c
}

// Unlike the ExtendByEnvs_WithPrefix(), this function allows to create new nodes in the config,
// based on the envs. It scans all envs matching the specified prefix, then strips the prefix,
// then what is left is used as a valid dot-path as is. For example, if the prefix was PRFX,
// and you specify an env var "PRFX_asd-qwe.zxc.123", then this variable will set, and create if
// necessary, the node at path "asd-qwe.zxc.123". If such names are supported in your OS is up
// to you, but see the https://stackoverflow.com/questions/2821043/allowed-characters-in-linux-environment-variable-names.
func (c *Config) ExtendByEnvsV2_WithPrefix(prefix string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			path := strings.Split(pair[0][len(prefix):], ".")
			value := pair[1]
			c.Set(path, value)
		}
	}
}
