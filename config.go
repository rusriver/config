package config

// Config ---------------------------------------------------------------------

// Config represents a configuration with convenient access methods.
type Config struct {
	Root    interface{}
	lastErr error
}

// Error return last error
func (c *Config) Error() error {
	return c.lastErr
}

// returns a nested config according to a dotted path.
func (c *Config) GetNestedConfig(path string) (*Config, error) {
	n, err := get(c.Root, path)
	if err != nil {
		return nil, err
	}
	return &Config{Root: n}, nil
}

// Sets a nested config according to a dotted path.
func (c *Config) Set(path string, v interface{}) error {
	return set(c.Root, path, v)
}

// Must is a wrapper for parsing functions to be used during initialization.
// It panics on failure.
func Must(c *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return c
}
