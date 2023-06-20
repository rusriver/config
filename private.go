package config

import (
	"fmt"
	"strconv"
	"strings"
)

func getAllPaths(source interface{}, base ...string) [][]string {
	paths := [][]string{}

	// Copy "base" so that underlying slice array is not
	// modified in recursive calls
	nextBase := make([]string, len(base))
	copy(nextBase, base)

	switch c := source.(type) {
	case map[string]interface{}:
		for k, v := range c {
			keys := getAllPaths(v, append(nextBase, k)...)
			paths = append(paths, keys...)
		}
	case []interface{}:
		for i, v := range c {
			k := strconv.Itoa(i)
			keys := getAllPaths(v, append(nextBase, k)...)
			paths = append(paths, keys...)
		}
	default:
		paths = append(paths, nextBase)
		return paths
	}
	return paths
}

// goByPath returns a child of the given value according to a dotted path.
func goByPath(c interface{}, pathParts []string) (interface{}, error) {
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

// Makes a path if necessary, (re-)sets a value
func set(c interface{}, pathParts []string, value interface{}) error {
	// Normalize the path.
	for k, v := range pathParts {
		if v == "" {
			if k == 0 {
				pathParts = pathParts[1:]
			} else {
				return fmt.Errorf("Invalid path %q", pathParts)
			}
		}
	}

	now := c
	for pathPart_i, pathPart_str := range pathParts {
		switch now_typed := now.(type) {

		case []interface{}:
			// we're in an array, check if the pp is a number, else err
			if i64, error := strconv.ParseInt(pathPart_str, 10, 0); error == nil {
				// it's a number, okay; we've also parsed it
				i := int(i64)

				// don't enlarge the array here, as it won't be saved in its parent then

				if pathPart_i+1 == len(pathParts) {
					// this is the last pp, which is indeed a number, we're in an array;
					// just set the value
					now_typed[i] = value
				} else {
					// this is NOT the last pp, there's more path

					if val := now_typed[i]; val != nil {
						// next thing exists just pick the pointer, and move in it

						// but first, make sure that IF the next pp is array, it is large enough,
						// and update the it HERE. Yeah, the it.
						switch val_typed := val.(type) {
						case []interface{}:
							next_pp := pathParts[pathPart_i+1]
							if i64, error := strconv.ParseInt(next_pp, 10, 0); error == nil {
								// it's a number, okay; we've also parsed it
								if int(i64) >= len(val_typed) {
									// enlarge array length if needed
									val_typed = append(val_typed, make([]interface{}, int(i64)-len(val_typed)+1)...)
								}
								// save it here
								now_typed[i] = val_typed
								val = val_typed
							}
						}

						// now move it it
						now = val
					} else {
						// we're in an array, and it doesn't have an element at index; we're about to create it;
						// is next path part a string (map key) or number (array index)?
						if i, err := strconv.ParseInt(pathParts[pathPart_i+1], 10, 0); err == nil {
							// next pp was a number, so create a nested array
							val = make([]interface{}, int(i)+1, int(i)+1)
						} else {
							// next pp was a string, to create a nested map
							val = make(map[string]interface{})
						}
						now_typed[i] = val
						now = val
					}
				}

			} else {
				return fmt.Errorf("Invalid list index at %q",
					strings.Join(pathParts[:pathPart_i+1], "."))
			}

		case map[string]interface{}:
			if pathPart_i+1 == len(pathParts) {
				// this is the last part of path, and we're in a map, therefore set value as map key value
				now_typed[pathPart_str] = value
			} else {
				// this is not the last pp
				if val, ok := now_typed[pathPart_str]; ok && val != nil {
					// we're in a map, it has our path part, so move in it, by pointer

					// but first, make sure that IF the next pp is array, it is large enough,
					// and update the it HERE. Yeah, the it.
					switch val_typed := val.(type) {
					case []interface{}:
						next_pp := pathParts[pathPart_i+1]
						if i64, error := strconv.ParseInt(next_pp, 10, 0); error == nil {
							// it's a number, okay; we've also parsed it
							if int(i64) >= len(val_typed) {
								// enlarge array length if needed
								val_typed = append(val_typed, make([]interface{}, int(i64)-len(val_typed)+1)...)
							}
							// save it here
							now_typed[pathPart_str] = val_typed
							val = val_typed
						}
					}

					// now move it it
					now = val
				} else {
					// we're in a map, and it doesn't have such key; we're about to create it;
					// is next path part a string (map key) or number (array index)?
					next_pp := pathParts[pathPart_i+1]
					if i, err := strconv.ParseInt(next_pp, 10, 0); err == nil {
						// next pp was a number, so create a nested array
						val = make([]interface{}, int(i)+1, int(i)+1)
					} else {
						// next pp was a string, to create a nested map
						val = make(map[string]interface{})
					}

					// set value as a map key value
					now_typed[pathPart_str] = val

					// move in further
					now = val
				}
			}

		default:
			return fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(pathParts[:pathPart_i+1], "."), now_typed)
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
