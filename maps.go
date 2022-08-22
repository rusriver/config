package config

import (
	"fmt"
	"strconv"
)

// MapConfig returns a nested map[string]*Config according to a dotted path.
func (c *Config) MapConfig(path string) (map[string]*Config, error) {
	m, err := c.Map(path)
	if err != nil {
		return nil, err
	}

	m2 := make(map[string]*Config, len(m))
	for k, v := range m {
		m2[k] = &Config{Root: v}
	}

	return m2, nil
}

// MapFloat64 returns a map[string]float64 according to a dotted path.
func (c *Config) MapFloat64(path string) (map[string]float64, error) {
	var err error

	m, err := c.Map(path)
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
		default:
			return m2, typeMismatch("float64, int or string", n)
		}
		m2[k] = v
	}

	return m2, nil
}

// MapInt returns a map[string]int according to a dotted path.
func (c *Config) MapInt(path string) (map[string]int, error) {
	var err error

	m, err := c.Map(path)
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
		default:
			return m2, typeMismatch("float64, int or string", n)
		}
		m2[k] = v
	}

	return m2, nil
}

// Map returns a map[string]interface{} according to a dotted path.
func (c *Config) Map(path string) (map[string]interface{}, error) {
	n, err := get(c.Root, path)
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

// MapString returns a map[string]string according to a dotted path.
func (c *Config) MapString(path string) (map[string]string, error) {
	var err error

	m, err := c.Map(path)
	if err != nil {
		return nil, err
	}

	m2 := make(map[string]string, len(m))
	for k, n := range m {
		var v string
		switch n := n.(type) {
		case string:
			v = n
		default:
			v = fmt.Sprintf("%v", n)
		}
		m2[k] = v
	}

	return m2, nil
}
