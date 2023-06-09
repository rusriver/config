package config

// Config ---------------------------------------------------------------------

// Config represents a configuration with convenient access methods.
type Config struct {
	Root          interface{}
	lastError     error
	ok            *bool
	err           *error
	dontPanicFlag bool
}

func (c *Config) Copy() *Config {
	return &Config{
		Root:          c.Root,
		ok:            c.ok,
		err:           c.err,
		dontPanicFlag: c.dontPanicFlag,
		lastError:     c.lastError,
	}
}

func (c *Config) LastError() error {
	return c.lastError
}

func (c *Config) P(path string) *Config {
	var err error
	c.Root, err = get(c.Root, path)
	if err != nil {
		c.handleError(err)
	}
	return c
}

func (c *Config) handleError(err error) {
	if err == nil {
		c.lastError = nil
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

// Sets a nested config according to a dotted path.
func (c *Config) Set(path string, v interface{}) {
	err := set(c.Root, path, v)
	if err != nil {
		c.handleError(err)
	}
}

func (c *Config) U() (c2 *Config) {
	c2 = c.Copy()
	c.dontPanicFlag = true
	return c2
}

func (c *Config) NotU() (c2 *Config) {
	c2 = c.Copy()
	c.dontPanicFlag = false
	return c2
}

func (c *Config) Ok(okRef *bool) (c2 *Config) {
	c2 = c.Copy()
	c.ok = okRef
	return c2
}

func (c *Config) E(err *error) (c2 *Config) {
	c2 = c.Copy()
	c.err = err
	return c2
}
