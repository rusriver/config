package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Get all keys for given interface, recursively
func getAllKeys(source interface{}, base ...string) [][]string {
	acc := [][]string{}

	// Copy "base" so that underlying slice array is not
	// modified in recursive calls
	nextBase := make([]string, len(base))
	copy(nextBase, base)

	switch c := source.(type) {
	case map[string]interface{}:
		for k, v := range c {
			keys := getAllKeys(v, append(nextBase, k)...)
			acc = append(acc, keys...)
		}
	case []interface{}:
		for i, v := range c {
			k := strconv.Itoa(i)
			keys := getAllKeys(v, append(nextBase, k)...)
			acc = append(acc, keys...)
		}
	default:
		acc = append(acc, nextBase)
		return acc
	}
	return acc
}

// get returns a child of the given value according to a dotted path.
func get(c interface{}, pathParts []string) (interface{}, error) {
	// Normalize path.
	for k, v := range pathParts {
		if v == "" {
			if k == 0 {
				pathParts = pathParts[1:]
			} else {
				return nil, fmt.Errorf("Invalid path %q", strings.Join(pathParts, "."))
			}
		}
	}
	// Get the value.
	for pos, part := range pathParts {
		switch cv := c.(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				if int(i) < len(cv) {
					c = cv[i]
				} else {
					return nil, fmt.Errorf(
						"Index out of range at %q: list has only %v items",
						strings.Join(pathParts[:pos+1], "."), len(cv))
				}
			} else {
				return nil, fmt.Errorf("Invalid list index at %q",
					strings.Join(pathParts[:pos+1], "."))
			}
		case map[string]interface{}:
			if value, ok := cv[part]; ok {
				c = value
			} else {
				return nil, fmt.Errorf("Nonexistent map key at %q",
					strings.Join(pathParts[:pos+1], "."))
			}
		default:
			return nil, fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(pathParts[:pos+1], "."), c)
		}
	}

	return c, nil
}

// set returns an error, in case when it is not possible to
// establish the value obtained in accordance with given dotted path.
func set(c interface{}, path string, value interface{}) error {
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

	point := &c
	for pos, part := range parts {
		switch c := (*point).(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				// 1. normalize slice capacity
				if int(i) >= len(c) {
					c = append(c, make([]interface{}, int(i)-len(c)+1)...)
					*point = c
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

// typeMismatchError returns an error for an expected type.
func typeMismatchError(expected string, got interface{}) error {
	return fmt.Errorf("Type mismatch: expected %s; got %T", expected, got)
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
