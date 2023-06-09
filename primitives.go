package config

import (
	"fmt"
	"strconv"
)

func (c *Config) Bool() bool {
	n := c.Root
	switch n := n.(type) {
	case bool:
		return n
	case string:
		b, err := strconv.ParseBool(n)
		if err != nil {
			c.handleError(err)
		}
		return b
	}
	c.handleError(typeMismatchError("bool or string", n))
	return false
}

func (c *Config) Float64() float64 {
	n := c.Root
	switch n := n.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case string:
		b, err := strconv.ParseFloat(n, 64)
		if err != nil {
			c.handleError(err)
		}
		return b
	}
	c.handleError(typeMismatchError("float64, int or string", n))
	return 0
}

func (c *Config) Int() int {
	n := c.Root
	switch n := n.(type) {
	case float64:
		// encoding/json unmarshals numbers into floats
		if i := int(n); float64(i) == n {
			return i
		} else {
			c.handleError(fmt.Errorf("Value can't be converted to int: %v", n))
		}
	case int:
		return n
	case string:
		if v, err := strconv.ParseInt(n, 10, 0); err == nil {
			return int(v)
		} else {
			c.handleError(err)
		}
	}
	c.handleError(typeMismatchError("float64, int or string", n))
	return 0
}

func (c *Config) String() string {
	n := c.Root
	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n)
	case string:
		return n
	}
	c.handleError(typeMismatchError("bool, float64, int or string", n))
	return ""
}
