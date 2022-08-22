package config

import (
	"fmt"
	"reflect"
)

// ExtendBy() extends current config with another config: i.e. all values
// from another config are added to the current config, and overwritten
// with new values if already present. It implements prototype-based inheritance.
func (c *Config) ExtendBy_v2(c2 *Config) (*Config, error) {

	newCfgInterface, err := extend(*c, *c2)
	if err != nil {
		return nil, err
	}

	newCfg, ok := newCfgInterface.(Config)
	if !ok {
		return nil, err
	}

	return &newCfg, nil
}

// Recursively extends c1 with c2
func extend(c1 interface{}, c2 interface{}) (interface{}, error) {

	// make sure types of ca and cb are correct

	if reflect.TypeOf(c1) != reflect.TypeOf(c2) {
		switch c1.(type) {
		case map[string]interface{}:
			c2 = map[string]interface{}{}
		case []interface{}:
			c2 = []interface{}{}
		case float64, string, int, bool:
			return c1, nil
			// do nothing
		default:
			return nil, fmt.Errorf("Invalid input. ca and cb must be of same type. They are %T %T\n", c1, c2)
		}
	}

	switch c1v := c1.(type) {
	case Config:
		c2v, _ := c2.(Config)
		c, err := extend(c1v.Root, c2v.Root)
		if err != nil {
			return nil, err
		}
		return Config{c, nil}, nil

	case map[string]interface{}:
		c2map, _ := c2.(map[string]interface{})
		c1map := c1v

		// Create a map and fill it with items from ca and cb
		newMap := map[string]interface{}{}
		var err error

		for k, v := range c1map {
			elementb, ok := c2map[k]
			if ok {
				newMap[k], err = extend(v, elementb)
				if err != nil {
					return newMap, err
				}
			} else {
				newMap[k], err = extend(v, nil)
				if err != nil {
					return newMap, err
				}
			}
		}

		for k, v := range c2map {
			_, ok := c1map[k]
			if !ok {
				val, err := extend(v, nil)
				if err != nil {
					return nil, err
				}
				newMap[k] = val
			}
		}
		return newMap, nil

	case []interface{}:
		c2arr, _ := c2.([]interface{})
		c1arr := c1v

		newArr := make([]interface{}, 0, cap(c1arr)+cap(c2arr))
		maxlen := len(c1arr)

		if len(c2arr) > len(c1arr) {
			maxlen = len(c2arr)
		}
		for i := 0; i < maxlen; i++ {

			if i < len(c1arr) {
				newArr = append(newArr, c1arr[i])
			} else if i < len(c2arr) {
				newArr = append(newArr, c2arr[i])
			} else {
				break
			}
		}
		return newArr, nil
	default:
		return c1, nil
	}
}
