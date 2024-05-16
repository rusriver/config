package config

import "time"

func (c *Config) Duration(defaultValueFunc ...func() time.Duration) time.Duration {
	n := c.DataSubTree
	if str, ok := n.(string); ok {
		dur, err := time.ParseDuration(str)
		if err == nil {
			return dur
		}
	}
	c.handleError(typeMismatchError("string", n))
	if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
		if c.ExpressionStatus == ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce {
			panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
		}
		c.ExpressionStatus++
		return defaultValueFunc[0]()
	} else {
		return 0
	}
}

func (c *Config) ListDuration(defaultValueFunc ...func() []time.Duration) []time.Duration {
	undef := make([]time.Duration, 0)
	l := c.List()

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
		c.handleError(typeMismatchError("string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionStatus == ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionStatus++
			return defaultValueFunc[0]()
		} else {
			return undef
		}
	OK:
		l2 = append(l2, v)
	}
	return l2
}

func (c *Config) MapDuration(defaultValueFunc ...func() map[string]time.Duration) map[string]time.Duration {
	undef := make(map[string]time.Duration)
	m := c.Map()

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
		c.handleError(typeMismatchError("string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionStatus == ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionStatus++
			return defaultValueFunc[0]()
		} else {
			return undef
		}
	OK:
		m2[k] = v
	}
	return m2
}
