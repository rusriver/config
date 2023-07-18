package config

import (
	"fmt"
	"strconv"
)

func (c *Config) Bool(defaultValueFunc ...func() bool) bool {
	n := c.DataSubTree
	switch n := n.(type) {
	case bool:
		return n
	case string:
		b, err := strconv.ParseBool(n)
		if err != nil {
			c.handleError(err)
			if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
				if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
					panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
				}
				c.ExpressionFailure++
				return defaultValueFunc[0]()
			} else {
				return false
			}
		}
		return b
	default:
		c.handleError(typeMismatchError("bool or string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return false
		}
	}
}

func (c *Config) Float64(defaultValueFunc ...func() float64) float64 {
	n := c.DataSubTree
	switch n := n.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case string:
		b, err := strconv.ParseFloat(n, 64)
		if err != nil {
			c.handleError(err)
			if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
				if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
					panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
				}
				c.ExpressionFailure++
				return defaultValueFunc[0]()
			} else {
				return 0
			}
		}
		return b
	default:
		c.handleError(typeMismatchError("float64, int or string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return 0
		}
	}
}

func (c *Config) Int(defaultValueFunc ...func() int) int {
	n := c.DataSubTree
	var err error
	switch n := n.(type) {
	case float64:
		// encoding/json unmarshals numbers into floats
		if i := int(n); float64(i) == n {
			return i
		}
		c.handleError(fmt.Errorf("Value can't be converted to int: %v", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return int(n)
		}
	case int:
		return n
	case string:
		if v, err := strconv.ParseInt(n, 10, 0); err == nil {
			return int(v)
		}
		c.handleError(err)
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return 0
		}
	default:
		c.handleError(typeMismatchError("float64, int or string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return 0
		}
	}
}

func (c *Config) String(defaultValueFunc ...func() string) string {
	n := c.DataSubTree
	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n)
	case string:
		return n
	default:
		c.handleError(typeMismatchError("bool, float64, int or string", n))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionFailure == ExpressionFailure_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionFailure++
			return defaultValueFunc[0]()
		} else {
			return fmt.Sprintf("%v", n)
		}
	}
}
