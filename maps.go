package config

import (
	"fmt"
	"strconv"
)

func (c *Config) Map() map[string]any {
	n := c.DataTreeRoot
	if value, ok := n.(map[string]interface{}); ok {
		return value
	}
	c.handleError(typeMismatchError("map[string]interface{}", n))
	return make(map[string]any)
}

func (c *Config) MapConfig() map[string]*Config {
	m := c.Map()

	m2 := make(map[string]*Config, len(m))
	for k, v := range m {
		m2[k] = c.Copy()
		m2[k].DataTreeRoot = v
	}

	return m2
}

func (c *Config) MapFloat64() map[string]float64 {
	m := c.Map()
	undef := make(map[string]float64)

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
				c.handleError(err)
				return undef
			}
			v = i
		default:
			c.handleError(typeMismatchError("float64, int or string", n))
			return undef
		}
		m2[k] = v
	}
	return m2
}

func (c *Config) MapInt() map[string]int {
	m := c.Map()
	undef := make(map[string]int)

	m2 := make(map[string]int, len(m))
	for k, n := range m {
		var v int
		switch n := n.(type) {
		case float64:
			// encoding/json unmarshals numbers into floats
			if i := int(n); float64(i) == n {
				v = i
			} else {
				c.handleError(fmt.Errorf("Value can't be converted to int: %v", n))
				return undef
			}
		case int:
			v = n
		case string:
			i, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				c.handleError(err)
				return undef
			}
			v = int(i)
		default:
			c.handleError(typeMismatchError("float64, int or string", n))
			return undef
		}
		m2[k] = v
	}
	return m2
}

func (c *Config) MapString() map[string]string {
	m := c.Map()

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

	return m2
}
