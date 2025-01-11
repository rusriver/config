package config

import "encoding/base64"

func (c *Config) BytesFromBase64(defaultValueFunc ...func() []byte) []byte {
	n := c.DataSubTree
	switch v := n.(type) {
	case string:
		bb, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			c.handleError(err)
			if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
				if c.ExpressionStatus == ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce {
					panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
				}
				c.ExpressionStatus++
				return defaultValueFunc[0]()
			} else {
				return []byte{}
			}
		}
		return bb
	default:
		c.handleError(typeMismatchError("string", v))
		if len(defaultValueFunc) > 0 && !c.isExpressionOk() {
			if c.ExpressionStatus == ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce {
				panic(ErrMsg_MultipleCallbackWithoutPriorErrOk)
			}
			c.ExpressionStatus++
			return defaultValueFunc[0]()
		} else {
			return []byte{}
		}
	}
}
