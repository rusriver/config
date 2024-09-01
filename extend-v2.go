package config

import (
	"reflect"

	"github.com/rusriver/config/v2/deepcopy"
)

// ExtendBy() extends current config with another config: i.e. all values
// from another config are added to the current config, and overwritten
// with new values if already present. It implements prototype-based inheritance.
func (c *Config) ExtendBy_v2(c2 *Config) *Config {
	c.DataSubTree = Extend_v2_any(c.DataSubTree, c2.DataSubTree)
	return c
}

// Recursively extends a1 with a2
func Extend_v2_any(a1 any, a2 any) any {

	if reflect.TypeOf(a1) == reflect.TypeOf(a2) {

		switch c1v := a1.(type) {

		case map[string]any:
			c2v := a2.(map[string]any)

			for k2, v2 := range c2v {
				switch v2v := v2.(type) {
				case map[string]any:
					if _, ok := c1v[k2]; !ok {
						c1v[k2] = make(map[string]any)
					}
					c1v[k2] = Extend_v2_any(c1v[k2], v2v)
				case []any:
					if _, ok := c1v[k2]; !ok {
						c1v[k2] = make([]any, 0)
					}
					c1v[k2] = Extend_v2_any(c1v[k2], v2v)
				default:
					c1v[k2] = v2
				}
			}

		case []any:
			c2v := a2.([]any)

			lenDiff := len(c2v) - len(c1v)
			if lenDiff > 0 {
				c1v = append(c1v, make([]any, lenDiff)...)
				a1 = c1v
			}

			for i2, v2 := range c2v {
				switch v2v := v2.(type) {
				case map[string]any:
					if c1v[i2] == nil {
						c1v[i2] = make(map[string]any)
					}
					c1v[i2] = Extend_v2_any(c1v[i2], v2v)
				case []any:
					if c1v[i2] == nil {
						c1v[i2] = make([]any, 0)
					}
					c1v[i2] = Extend_v2_any(c1v[i2], v2v)
				default:
					c1v[i2] = v2
				}
			}

		}

		// END if reflect.TypeOf(c1) == reflect.TypeOf(c2)
	} else {
		a2_copy := deepcopy.Copy(a2)
		return a2_copy
	}

	return a1
}
