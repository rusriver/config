package config

import (
	"bytes"
)

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

func (c *Config) P(pathParts ...string) *Config {
	c.resetErrOkState()
	c = c.Copy()
	var err error
	c.Root, err = get(c.Root, pathParts)
	if err != nil {
		c.handleError(err)
	}
	return c
}

func (c *Config) resetErrOkState() {
	c.lastError = nil
	if c.ok != nil {
		*c.ok = true
	}
	if c.err != nil {
		*c.err = nil
	}

}

func (c *Config) handleError(err error) {
	if err == nil {
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
func (c *Config) Set(pathParts []string, v interface{}) {
	c.resetErrOkState()
	err := set(c.Root, pathParts, v)
	if err != nil {
		c.handleError(err)
	}
}

func (c *Config) U() (c2 *Config) {
	c.resetErrOkState()
	c2 = c.Copy()
	c2.dontPanicFlag = true
	return c2
}

func (c *Config) NotU() (c2 *Config) {
	c.resetErrOkState()
	c2 = c.Copy()
	c2.dontPanicFlag = false
	return c2
}

func (c *Config) Ok(okRef *bool) (c2 *Config) {
	c.resetErrOkState()
	c2 = c.Copy()
	c2.ok = okRef
	return c2
}

func (c *Config) E(err *error) (c2 *Config) {
	c.resetErrOkState()
	c2 = c.Copy()
	c2.err = err
	return c2
}

func SplitPathToParts(key string) []string {
	parts := []string{}

	bracketOpened := false
	var buffer bytes.Buffer
	for _, char := range key {
		if char == 91 || char == 93 { // [ ]
			bracketOpened = char == 91
			continue
		}
		if char == 46 && !bracketOpened { // point
			parts = append(parts, buffer.String())
			buffer.Reset()
			continue
		}

		buffer.WriteRune(char)
	}

	if buffer.String() != "" {
		parts = append(parts, buffer.String())
		buffer.Reset()
	}

	return parts
}
