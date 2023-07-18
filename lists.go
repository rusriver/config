package config

import (
	"fmt"
	"strconv"
)

func (c *Config) List(defaultValueFunc ...func() []any) []any {
	n := c.DataSubTree
	if value, ok := n.([]interface{}); ok {
		return value
	}
	c.handleError(typeMismatchError("[]interface{}", n))
	if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
		if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
			panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
		}
		c.ExpressionFailure++
		return defaultValueFunc[0]()
	} else {
		return make([]any, 0)
	}
}

func (c *Config) ListConfig() []*Config {
	l := c.List()

	l2 := make([]*Config, 0, len(l))
	for _, v := range l {
		c2 := c.ChildCopy()
		c2.DataSubTree = v
		l2 = append(l2, c2)
	}

	return l2
}

func (c *Config) ListFloat64(defaultValueFunc ...func() []float64) []float64 {
	l := c.List()
	undef := make([]float64, 0)

	l2 := make([]float64, 0, len(l))
	for _, n := range l {
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
				if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
					if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
						panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
					}
					c.ExpressionFailure++
					return defaultValueFunc[0]()
				} else {
					return undef
				}
			}
			v = i
		default:
			c.handleError(typeMismatchError("float64, int or string", n))
			if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
				if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
					panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
				}
				c.ExpressionFailure++
				return defaultValueFunc[0]()
			} else {
				return undef
			}
		}
		l2 = append(l2, v)
	}
	return l2
}

func (c *Config) ListInt(defaultValueFunc ...func() []int) []int {
	l := c.List()
	undef := make([]int, 0)

	l2 := make([]int, 0, len(l))
	for _, n := range l {
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
				if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
					if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
						panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
					}
					c.ExpressionFailure++
					return defaultValueFunc[0]()
				} else {
					return undef
				}
			}
			v = int(i)
		default:
			c.handleError(typeMismatchError("float64, int or string", n))
			if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
				if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
					panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
				}
				c.ExpressionFailure++
				return defaultValueFunc[0]()
			} else {
				return undef
			}
		}
		l2 = append(l2, v)
	}
	return l2
}

func (c *Config) ListString(defaultValueFunc ...func() []string) []string {
	if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
		if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
			panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
		}
		c.ExpressionFailure++
		return defaultValueFunc[0]()
	}

	l := c.List()

	l2 := make([]string, 0, len(l))
	for _, n := range l {
		var v string
		switch n := n.(type) {
		case string:
			v = n
		default:
			v = fmt.Sprintf("%v", n)
		}
		l2 = append(l2, v)
	}

	return l2
}
