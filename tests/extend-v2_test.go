package main

import (
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_Extend_v2_1_0(t *testing.T) {
	c1 := (&config.InitContext{}).
		FromFile("conf-test-files/extend-v2/20240901-1.yaml").Load()

	c2 := (&config.InitContext{}).
		FromFile("conf-test-files/extend-v2/20240901-2.yaml").Load()

	c1.ExtendBy_v2(c2)
	c1.PrintJson("RESULT")

	c3 := (&config.InitContext{}).
		FromFile("conf-test-files/extend-v2/20240901-3.yaml").Load()

	c1.ExtendBy_v2(c3)
	c1.PrintJson("RESULT")
}
