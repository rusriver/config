package config

import "time"

// Duration returns a time.Duration according to a dotted path.
func (c *Config) Duration(path string) (time.Duration, error) {
	n, err := get(c.Root, path)
	if err != nil {
		return 0, err
	}
	if str, ok := n.(string); ok {
		dur, err := time.ParseDuration(str)
		if err == nil {
			return dur, nil
		}
	}
	return 0, typeMismatch("string", n)
}

// UDuration returns a time.Duration according to a dotted path or default or zero duration
func (c *Config) UDuration(path string, defaults ...time.Duration) time.Duration {
	dur, err := c.Duration(path)

	if err == nil {
		return dur
	}

	for _, def := range defaults {
		return def
	}
	return 0
}

// ListDuration returns an []time.Duration according to a dotted path.
func (c *Config) ListDuration(path string) ([]time.Duration, error) {
	l, err := c.List(path)
	if err != nil {
		return nil, err
	}

	l2 := make([]time.Duration, 0, len(l))
	for _, n := range l {
		var v time.Duration
		if str, ok := n.(string); ok {
			dur, err := time.ParseDuration(str)
			if err == nil {
				v = dur
				goto OK
			}
		}
		return l2, typeMismatch("string", n)
	OK:
		l2 = append(l2, v)
	}

	return l2, nil
}

// MapDuration returns a map[string]time.Duration according to a dotted path.
func (c *Config) MapDuration(path string) (map[string]time.Duration, error) {
	var err error

	m, err := c.Map(path)
	if err != nil {
		return nil, err
	}

	m2 := make(map[string]time.Duration, len(m))
	for k, n := range m {
		var v time.Duration
		if str, ok := n.(string); ok {
			dur, err := time.ParseDuration(str)
			if err == nil {
				v = dur
				goto OK
			}
		}
		return m2, typeMismatch("string", n)
	OK:
		m2[k] = v
	}

	return m2, nil
}
