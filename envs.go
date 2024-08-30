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
// Update: As it turns out, the OS may support it, but the bash didn't, and so you can use it.
// To solve this issue, I am adding a transformation callback, which can additionally transform
// the path string, so you can avoid using non-alphanumerics in env names.
// See also the TranslateEnvs_KeySuffix().
func (c *Config) ExtendByEnvsV2_WithPrefix(prefix string, transformerFuncs ...func(s string) string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			for _, tf := range transformerFuncs {
				pair[0] = tf(pair[0])
			}
			path := strings.Split(pair[0][len(prefix):], ".")
			value := pair[1]
			c.Set(path, value)
		}
	}
}

// Does this replacements:
//
//	p__ -> "." (point)
//	d__ -> "-" (dash)
//	u__ -> "__" (underscore)
//	D__ -> "$" (dollar)
//
// For example, the line "Somethingp__superd__duperp__D__isa" will become "Something.super-duper.$isa".
// Because the key is located at the end of words, this makes minimal possible impact on readability.
func TranslateEnvs_KeySuffix(s string) string {
	s = strings.ReplaceAll(s, "p__", ".")
	s = strings.ReplaceAll(s, "d__", "-")
	s = strings.ReplaceAll(s, "u__", "__")
	return s
}
