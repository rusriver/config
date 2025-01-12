package config

import (
	"time"

	"github.com/rusriver/dateparse"
)

// Parses almost any date-time format
func (c *Config) Time(defaultValueFunc ...func() time.Time) time.Time {
	n := c.DataSubTree
	if str, ok := n.(string); ok {
		t, err := dateparse.ParseStrict(str)
		if err == nil {
			return t
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
		return time.Unix(0, 0)
	}
}

func (c *Config) ListTime(defaultValueFunc ...func() []time.Time) []time.Time {
	undef := make([]time.Time, 0)
	l := c.List()

	l2 := make([]time.Time, 0, len(l))
	for _, n := range l {
		var v time.Time
		if str, ok := n.(string); ok {
			t, err := dateparse.ParseStrict(str)
			if err == nil {
				v = t
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

func (c *Config) MapTime(defaultValueFunc ...func() map[string]time.Time) map[string]time.Time {
	undef := make(map[string]time.Time)
	m := c.Map()

	m2 := make(map[string]time.Time, len(m))
	for k, n := range m {
		var v time.Time
		if str, ok := n.(string); ok {
			t, err := dateparse.ParseStrict(str)
			if err == nil {
				v = t
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
