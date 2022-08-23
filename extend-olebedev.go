package config

import "strings"

// ExtendBy() extends current config with another config: i.e. all values
// from another config are added to the current config, and overwritten
// with new values if already present. It implements prototype-based inheritance.
// Note that if you extend with different structure you will get an error.
// See: `.Set()` method for details.
func (c *Config) ExtendBy(c2 *Config) (err error) {
	keys := getKeys(c2.Root)
	for _, key := range keys {
		k := strings.Join(key, ".")
		i, err := get(c2.Root, k)
		if err != nil {
			return err
		}
		if err := c.Set(k, i); err != nil {
			return err
		}
	}
	return nil
}