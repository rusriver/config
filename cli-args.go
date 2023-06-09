package config

import (
	"bytes"
	"flag"
	"strings"
)

// Parse command line arguments, based on existing config keys.
func (c *Config) Flag() *Config {
	keys := getAllKeys(c.Root)
	hash := map[string]*string{}
	for _, key := range keys {
		k := strings.Join(key, "-")
		hash[k] = new(string)
		val := c.P(k).String()
		flag.StringVar(hash[k], k, val, "")
	}

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		name := strings.Replace(f.Name, "-", ".", -1)
		c.Set(name, f.Value.String())
	})

	return c
}

// Args command line arguments, based on existing config keys.
func (c *Config) Args(args ...string) *Config {
	if len(args) <= 1 {
		return c
	}

	keys := getAllKeys(c.Root)
	hash := map[string]*string{}
	_flag := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var _err bytes.Buffer
	_flag.SetOutput(&_err)
	for _, key := range keys {
		k := strings.Join(key, "-")
		hash[k] = new(string)
		val := c.P(k).String()
		_flag.StringVar(hash[k], k, val, "")
	}

	c.lastError = _flag.Parse(args[1:])

	_flag.Visit(func(f *flag.Flag) {
		name := strings.Replace(f.Name, "-", ".", -1)
		c.Set(name, f.Value.String())
	})

	return c
}
