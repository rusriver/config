// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	yaml "gopkg.in/yaml.v2"
)

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

// DEPRECATED, use Config() instead.
func (cfg *Config) Get(path string) (*Config, error) {
	return nil, nil
}

// Config returns a nested config according to a dotted path.
func (cfg *Config) Config(path string) (*Config, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	return &Config{Root: n}, nil
}

// ListConfig returns a nested []*Config according to a dotted path.
func (cfg *Config) ListConfig(path string) ([]*Config, error) {
	l, err := cfg.List(path)
	if err != nil {
		return nil, err
	}
	
	l2 := make([]*Config, 0, len(l))
	for _, v := range l {
		l2 = append(l2, &Config{Root: v})
	}
	
	return l2, nil
}

// MapConfig returns a nested map[string]*Config according to a dotted path.
func (cfg *Config) MapConfig(path string) (map[string]*Config, error) {
	m, err := cfg.Map(path)
	if err != nil {
		return nil, err
	}
	
	m2 := make(map[string]*Config, len(m))
	for k, v := range m {
		m2[k] = &Config{Root: v}
	}
	
	return m2, nil
}

// Set a nested config according to a dotted path.
func (cfg *Config) Set(path string, val interface{}) error {
	return Set(cfg.Root, path, val)
}

// Fetch data from system env, based on existing config keys.
func (cfg *Config) Env() *Config {
	return cfg.EnvPrefix("")
}

// Fetch data from system env using prefix, based on existing config keys.
func (cfg *Config) EnvPrefix(prefix string) *Config {
	if prefix != "" {
		prefix = strings.ToUpper(prefix) + "_"
	}

	keys := getKeys(cfg.Root)
	for _, key := range keys {
		k := strings.ToUpper(strings.Join(key, "_"))
		if val, exist := syscall.Getenv(prefix + k); exist {
			cfg.Set(strings.Join(key, "."), val)
		}
	}
	return cfg
}

// Parse command line arguments, based on existing config keys.
func (cfg *Config) Flag() *Config {
	keys := getKeys(cfg.Root)
	hash := map[string]*string{}
	for _, key := range keys {
		k := strings.Join(key, "-")
		hash[k] = new(string)
		val, _ := cfg.String(k)
		flag.StringVar(hash[k], k, val, "")
	}

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		name := strings.Replace(f.Name, "-", ".", -1)
		cfg.Set(name, f.Value.String())
	})

	return cfg
}

// Args command line arguments, based on existing config keys.
func (cfg *Config) Args(args ...string) *Config {
	if len(args) <= 1 {
		return cfg
	}

	keys := getKeys(cfg.Root)
	hash := map[string]*string{}
	_flag := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var _err bytes.Buffer
	_flag.SetOutput(&_err)
	for _, key := range keys {
		k := strings.Join(key, "-")
		hash[k] = new(string)
		val, _ := cfg.String(k)
		_flag.StringVar(hash[k], k, val, "")
	}

	cfg.lastErr = _flag.Parse(args[1:])

	_flag.Visit(func(f *flag.Flag) {
		name := strings.Replace(f.Name, "-", ".", -1)
		cfg.Set(name, f.Value.String())
	})

	return cfg
}

// Get all keys for given interface
func getKeys(source interface{}, base ...string) [][]string {
	acc := [][]string{}

	// Copy "base" so that underlying slice array is not
	// modified in recursive calls
	nextBase := make([]string, len(base))
	copy(nextBase, base)

	switch c := source.(type) {
	case map[string]interface{}:
		for k, v := range c {
			keys := getKeys(v, append(nextBase, k)...)
			acc = append(acc, keys...)
		}
	case []interface{}:
		for i, v := range c {
			k := strconv.Itoa(i)
			keys := getKeys(v, append(nextBase, k)...)
			acc = append(acc, keys...)
		}
	default:
		acc = append(acc, nextBase)
		return acc
	}
	return acc
}

// Bool returns a bool according to a dotted path.
func (cfg *Config) Bool(path string) (bool, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return false, err
	}
	switch n := n.(type) {
	case bool:
		return n, nil
	case string:
		return strconv.ParseBool(n)
	}
	return false, typeMismatch("bool or string", n)
}

// UBool returns a bool according to a dotted path or default value or false.
func (c *Config) UBool(path string, defaults ...bool) bool {
	value, err := c.Bool(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return false
}

// Float64 returns a float64 according to a dotted path.
func (cfg *Config) Float64(path string) (float64, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return 0, err
	}
	switch n := n.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case string:
		return strconv.ParseFloat(n, 64)
	}
	return 0, typeMismatch("float64, int or string", n)
}

// UFloat64 returns a float64 according to a dotted path or default value or 0.
func (c *Config) UFloat64(path string, defaults ...float64) float64 {
	value, err := c.Float64(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return float64(0)
}

// ListFloat64 returns an []float64 according to a dotted path.
func (cfg *Config) ListFloat64(path string) ([]float64, error) {
	l, err := cfg.List(path)
	if err != nil {
		return nil, err
	}
	
	l2 := make([]float64, 0, len(l))
	for _, n := range l {
		var v float64
		switch n := n.(type) {
		case float64:
			// encoding/json unmarshals numbers into floats
			v = n
		case int:
			v = float64(n)
		case string:
			i, err := strconv.ParseFloat(n, 64)
			if err != nil {
				return l2, err
			}
			v = i
		}
		l2 = append(l2, v)
	}
	
	return l2, nil
}

// MapFloat64 returns a map[string]float64 according to a dotted path.
func (cfg *Config) MapFloat64(path string) (map[string]float64, error) {
	var err error
	
	m, err := cfg.Map(path)
	if err != nil {
		return nil, err
	}
	
	m2 := make(map[string]float64, len(m))
	for k, n := range m {
		var v float64
		switch n := n.(type) {
		case float64:
			// encoding/json unmarshals numbers into floats
			v = n
		case int:
			v = float64(n)
		case string:
			i, err := strconv.ParseFloat(n, 64)
			if err != nil {
				return m2, err
			}
			v = i
		}
		m2[k] = v
	}
	
	return m2, nil
}

// Int returns an int according to a dotted path.
func (cfg *Config) Int(path string) (int, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return 0, err
	}
	switch n := n.(type) {
	case float64:
		// encoding/json unmarshals numbers into floats
		if i := int(n); float64(i) == n {
			return i, nil
		} else {
			return 0, fmt.Errorf("Value can't be converted to int: %v", n)
		}
	case int:
		return n, nil
	case string:
		if v, err := strconv.ParseInt(n, 10, 0); err == nil {
			return int(v), nil
		} else {
			return 0, err
		}
	}
	return 0, typeMismatch("float64, int or string", n)
}

// UInt returns an int according to a dotted path or default value or 0.
func (c *Config) UInt(path string, defaults ...int) int {
	value, err := c.Int(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return 0
}

// ListInt returns an []int according to a dotted path.
func (cfg *Config) ListInt(path string) ([]int, error) {
	l, err := cfg.List(path)
	if err != nil {
		return nil, err
	}
	
	l2 := make([]int, 0, len(l))
	for _, n := range l {
		var v int
		switch n := n.(type) {
		case float64:
			// encoding/json unmarshals numbers into floats
			if i := int(n); float64(i) == n {
				v = i
			} else {
				return l2, fmt.Errorf("Value can't be converted to int: %v", n)
			}
		case int:
			v = n
		case string:
			i, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				return l2, err
			}
			v = int(i)
		}
		l2 = append(l2, v)
	}
	
	return l2, nil
}

// MapInt returns a map[string]int according to a dotted path.
func (cfg *Config) MapInt(path string) (map[string]int, error) {
	var err error
	
	m, err := cfg.Map(path)
	if err != nil {
		return nil, err
	}
	
	m2 := make(map[string]int, len(m))
	for k, n := range m {
		var v int
		switch n := n.(type) {
		case float64:
			// encoding/json unmarshals numbers into floats
			if i := int(n); float64(i) == n {
				v = i
			} else {
				return m2, fmt.Errorf("Value can't be converted to int: %v", n)
			}
		case int:
			v = n
		case string:
			i, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				return m2, err
			}
			v = int(i)
		}
		m2[k] = v
	}
	
	return m2, nil
}

// List returns a []interface{} according to a dotted path.
func (cfg *Config) List(path string) ([]interface{}, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	if value, ok := n.([]interface{}); ok {
		return value, nil
	}
	return nil, typeMismatch("[]interface{}", n)
}

// UList returns a []interface{} according to a dotted path or defaults or []interface{}.
func (c *Config) UList(path string, defaults ...[]interface{}) []interface{} {
	value, err := c.List(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return make([]interface{}, 0)
}

// Map returns a map[string]interface{} according to a dotted path.
func (cfg *Config) Map(path string) (map[string]interface{}, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	if value, ok := n.(map[string]interface{}); ok {
		return value, nil
	}
	return nil, typeMismatch("map[string]interface{}", n)
}

// UMap returns a map[string]interface{} according to a dotted path or default or map[string]interface{}.
func (c *Config) UMap(path string, defaults ...map[string]interface{}) map[string]interface{} {
	value, err := c.Map(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return map[string]interface{}{}
}

// String returns a string according to a dotted path.
func (cfg *Config) String(path string) (string, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return "", err
	}
	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n), nil
	case string:
		return n, nil
	}
	return "", typeMismatch("bool, float64, int or string", n)
}

// UString returns a string according to a dotted path or default or "".
func (c *Config) UString(path string, defaults ...string) string {
	value, err := c.String(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return ""
}

// Duration returns a time.Duration according to a dotted path.
func (cfg *Config) Duration(path string) (time.Duration, error) {
	n, err := get(cfg.Root, path)
	if err != nil {
		return 0, err
	}
	if str, ok := n.(string); ok {
		dur, err := time.ParseDuration(str)
		if err == nil {
			return dur, nil
		}
	}
	return 0, typeMismatch("string", n)
}

// UDuration returns a time.Duration according to a dotted path or default or zero duration
func (cfg *Config) UDuration(path string, defaults ...time.Duration) (time.Duration) {
	dur, err := cfg.Duration(path)
	
	if err == nil {
		return dur
	}

	for _, def := range defaults {
		return def
	}
	return 0
}

// Copy returns a deep copy with given path or without.
func (c *Config) Copy(dottedPath ...string) (*Config, error) {
	toJoin := []string{}
	for _, part := range dottedPath {
		if len(part) != 0 {
			toJoin = append(toJoin, part)
		}
	}

	var err error
	var path = strings.Join(toJoin, ".")
	var cfg = c
	var root = ""

	if len(path) > 0 {
		if cfg, err = c.Config(path); err != nil {
			return nil, err
		}
	}

	if root, err = RenderYaml(cfg.Root); err != nil {
		return nil, err
	}
	return ParseYaml(root)
}

// Extend returns extended copy of current config with applied
// values from the given config instance. Note that if you extend
// with different structure you will get an error. See: `.Set()` method
// for details.
func (c *Config) Extend(cfg *Config) (*Config, error) {
	n, err := c.Copy()
	if err != nil {
		return nil, err
	}

	keys := getKeys(cfg.Root)
	for _, key := range keys {
		k := strings.Join(key, ".")
		i, err := get(cfg.Root, k)
		if err != nil {
			return nil, err
		}
		if err := n.Set(k, i); err != nil {
			return nil, err
		}
	}
	return n, nil
}

// typeMismatch returns an error for an expected type.
func typeMismatch(expected string, got interface{}) error {
	return fmt.Errorf("Type mismatch: expected %s; got %T", expected, got)
}

// Fetching -------------------------------------------------------------------

// get returns a child of the given value according to a dotted path.
func get(cfg interface{}, path string) (interface{}, error) {
	parts := splitKeyOnParts(path)
	// Normalize path.
	for k, v := range parts {
		if v == "" {
			if k == 0 {
				parts = parts[1:]
			} else {
				return nil, fmt.Errorf("Invalid path %q", path)
			}
		}
	}
	// Get the value.
	for pos, part := range parts {
		switch c := cfg.(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				if int(i) < len(c) {
					cfg = c[i]
				} else {
					return nil, fmt.Errorf(
						"Index out of range at %q: list has only %v items",
						strings.Join(parts[:pos+1], "."), len(c))
				}
			} else {
				return nil, fmt.Errorf("Invalid list index at %q",
					strings.Join(parts[:pos+1], "."))
			}
		case map[string]interface{}:
			if value, ok := c[part]; ok {
				cfg = value
			} else {
				return nil, fmt.Errorf("Nonexistent map key at %q",
					strings.Join(parts[:pos+1], "."))
			}
		default:
			return nil, fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(parts[:pos+1], "."), cfg)
		}
	}

	return cfg, nil
}

func splitKeyOnParts(key string) []string {
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

// Set returns an error, in case when it is not possible to
// establish the value obtained in accordance with given dotted path.
func Set(cfg interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	// Normalize path.
	for k, v := range parts {
		if v == "" {
			if k == 0 {
				parts = parts[1:]
			} else {
				return fmt.Errorf("Invalid path %q", path)
			}
		}
	}

	point := &cfg
	for pos, part := range parts {
		switch c := (*point).(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				// 1. normalize slice capacity
				if int(i) >= cap(c) {
					c = append(c, make([]interface{}, int(i)-cap(c)+1, int(i)-cap(c)+1)...)
				}

				// 2. set value or go further
				if pos+1 == len(parts) {
					c[i] = value
				} else {

					// if exists just pick the pointer
					if va := c[i]; va != nil {
						point = &va
					} else {
						// is next part slice or map?
						if i, err := strconv.ParseInt(parts[pos+1], 10, 0); err == nil {
							va = make([]interface{}, int(i)+1, int(i)+1)
						} else {
							va = make(map[string]interface{})
						}
						c[i] = va
						point = &va
					}

				}

			} else {
				return fmt.Errorf("Invalid list index at %q",
					strings.Join(parts[:pos+1], "."))
			}
		case map[string]interface{}:
			if pos+1 == len(parts) {
				c[part] = value
			} else {
				// if exists just pick the pointer
				if va, ok := c[part]; ok && va != nil {
					point = &va
				} else {
					// is next part slice or map?
					if i, err := strconv.ParseInt(parts[pos+1], 10, 0); err == nil {
						va = make([]interface{}, int(i)+1, int(i)+1)
					} else {
						va = make(map[string]interface{})
					}
					c[part] = va
					point = &va
				}
			}
		default:
			return fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(parts[:pos+1], "."), c)
		}
	}

	return nil
}

// Parsing --------------------------------------------------------------------

// Must is a wrapper for parsing functions to be used during initialization.
// It panics on failure.
func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

// normalizeValue normalizes a unmarshalled value. This is needed because
// encoding/json doesn't support marshalling map[interface{}]interface{}.
func normalizeValue(value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case map[interface{}]interface{}:
		node := make(map[string]interface{}, len(value))
		for k, v := range value {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Unsupported map key: %#v", k)
			}
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case map[string]interface{}:
		node := make(map[string]interface{}, len(value))
		for key, v := range value {
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case []interface{}:
		node := make([]interface{}, len(value))
		for key, v := range value {
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported list item: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case bool, float64, int, string, nil:
		return value, nil
	}
	return nil, fmt.Errorf("Unsupported type: %T", value)
}

// JSON -----------------------------------------------------------------------

// ParseJson reads a JSON configuration from the given string.
func ParseJson(cfg string) (*Config, error) {
	return parseJson([]byte(cfg))
}

// ParseJsonFile reads a JSON configuration from the given filename.
func ParseJsonFile(filename string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseJson(cfg)
}

// parseJson performs the real JSON parsing.
func parseJson(cfg []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = json.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderJson renders a JSON configuration.
func RenderJson(cfg interface{}) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// YAML -----------------------------------------------------------------------

// ParseYamlBytes reads a YAML configuration from the given []byte.
func ParseYamlBytes(cfg []byte) (*Config, error) {
	return parseYaml(cfg)
}

// ParseYaml reads a YAML configuration from the given string.
func ParseYaml(cfg string) (*Config, error) {
	return parseYaml([]byte(cfg))
}

// ParseYamlFile reads a YAML configuration from the given filename.
func ParseYamlFile(filename string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseYaml(cfg)
}

// parseYaml performs the real YAML parsing.
func parseYaml(cfg []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = yaml.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderYaml renders a YAML configuration.
func RenderYaml(cfg interface{}) (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
