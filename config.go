package config

// Config ---------------------------------------------------------------------

// Config represents a configuration with convenient access methods.
type Config struct {
	Root            interface{}
	lastError       error
	ok              *bool
	err             *error
	dontPanicFlag   bool
	lastFetchedNode any
}

// LastError return last error
func (c *Config) LastError() error {
	return c.lastError
}

func (c *Config) P(path string) *Config {
	var err error
	c.lastFetchedNode, err = get(c.Root, path)
	if err != nil {
		c.handleError(err)
	}
	return c
}

func (c *Config) handleError(err error) {
	if err == nil {
		if c.ok != nil {
			*c.ok = true
		}
		return
	} else {
		c.lastError = err
		if c.err != nil {
			*c.err = err
		}
		if c.ok != nil {
			*c.ok = false
		}
		if c.err == nil && c.ok == nil && !c.dontPanicFlag {
			panic(err)
		}
	}
}

func (c *Config) NestedConfig() *Config {
	return &Config{Root: c.lastFetchedNode} // don't copy ok, err, u
}

// Sets a nested config according to a dotted path.
func (c *Config) Set(path string, v interface{}) {
	err := set(c.Root, path, v)
	if err != nil {
		c.handleError(err)
	}
}

func (c *Config) U() (c2 *Config) {
	c2 = &Config{
		Root:          c.Root,
		ok:            c.ok,
		err:           c.err,
		dontPanicFlag: true,
	}
	return c2
}

func (c *Config) Ok(okRef *bool) (c2 *Config) {
	c2 = &Config{
		Root:          c.Root,
		ok:            okRef,
		err:           c.err,
		dontPanicFlag: c.dontPanicFlag,
	}
	return c2
}

func (c *Config) E(err *error) (c2 *Config) {
	c2 = &Config{
		Root:          c.Root,
		ok:            c.ok,
		err:           err,
		dontPanicFlag: c.dontPanicFlag,
	}
	return c2
}
