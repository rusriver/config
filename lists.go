package config

import (
	"fmt"
	"strconv"
)

// ListConfig returns a nested []*Config according to a dotted path.
func (c *Config) ListConfig(path string) ([]*Config, error) {
	l, err := c.List(path)
	if err != nil {
		return nil, err
	}

	l2 := make([]*Config, 0, len(l))
	for _, v := range l {
		l2 = append(l2, &Config{Root: v})
	}

	return l2, nil
}

// ListFloat64 returns an []float64 according to a dotted path.
func (c *Config) ListFloat64(path string) ([]float64, error) {
	l, err := c.List(path)
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
		default:
			return l2, typeMismatch("float64, int or string", n)
		}
		l2 = append(l2, v)
	}

	return l2, nil
}

// ListInt returns an []int according to a dotted path.
func (c *Config) ListInt(path string) ([]int, error) {
	l, err := c.List(path)
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
		default:
			return l2, typeMismatch("float64, int or string", n)
		}
		l2 = append(l2, v)
	}

	return l2, nil
}

// List returns a []interface{} according to a dotted path.
func (c *Config) List(path string) ([]interface{}, error) {
	n, err := get(c.Root, path)
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

// ListString returns an []string according to a dotted path.
func (c *Config) ListString(path string) ([]string, error) {
	l, err := c.List(path)
	if err != nil {
		return nil, err
	}

	l2 := make([]string, 0, len(l))
	for _, n := range l {
		var v string
		switch n := n.(type) {
		case string:
			v = n
		default:
			v = fmt.Sprintf("%v", n)
		}
		l2 = append(l2, v)
	}

	return l2, nil
}
