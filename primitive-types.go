package config

import (
	"fmt"
	"strconv"
)

// Bool returns a bool according to a dotted path.
func (c *Config) Bool(path string) (bool, error) {
	n, err := get(c.Root, path)
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
func (c *Config) Float64(path string) (float64, error) {
	n, err := get(c.Root, path)
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

// Int returns an int according to a dotted path.
func (c *Config) Int(path string) (int, error) {
	n, err := get(c.Root, path)
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

// String returns a string according to a dotted path.
func (c *Config) String(path string) (string, error) {
	n, err := get(c.Root, path)
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
