package config

import (
	"reflect"
)

// ExtendBy() extends current config with another config: i.e. all values
// from another config are added to the current config, and overwritten
// with new values if already present. It implements prototype-based inheritance.
func (c *Config) ExtendBy_v2(c2 *Config) *Config {
	c.DataTreeRoot = extend_v2(c.DataTreeRoot, c2.DataTreeRoot)
	return c
}

// Recursively extends c1 with c2
func extend_v2(c1 interface{}, c2 interface{}) interface{} {

	if reflect.TypeOf(c1) == reflect.TypeOf(c2) {

		switch c1v := c1.(type) {

		case map[string]interface{}:
			c2v := c2.(map[string]interface{})

			for k2, v2 := range c2v {
				switch v2v := v2.(type) {
				case map[string]interface{}:
					if _, ok := c1v[k2]; !ok {
						c1v[k2] = make(map[string]interface{})
					}
					c1v[k2] = extend_v2(c1v[k2], v2v)
				case []interface{}:
					if _, ok := c1v[k2]; !ok {
						c1v[k2] = make([]interface{}, 0)
					}
					c1v[k2] = extend_v2(c1v[k2], v2v)
				default:
					c1v[k2] = v2
				}
			}

		case []interface{}:
			c2v := c2.([]interface{})

			lenDiff := len(c2v) - len(c1v)
			if lenDiff > 0 {
				c1v = append(c1v, make([]interface{}, lenDiff)...)
				c1 = c1v
			}

			for i2, v2 := range c2v {
				switch v2v := v2.(type) {
				case map[string]interface{}:
					if c1v[i2] == nil {
						c1v[i2] = make(map[string]any)
					}
					c1v[i2] = extend_v2(c1v[i2], v2v)
				case []interface{}:
					if c1v[i2] == nil {
						c1v[i2] = make([]any, 0)
					}
					c1v[i2] = extend_v2(c1v[i2], v2v)
				default:
					c1v[i2] = v2
				}
			}

		}

	} // if reflect.TypeOf(c1) == reflect.TypeOf(c2)

	return c1
}
