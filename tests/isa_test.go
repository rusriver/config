package main

import (
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_Isa_1_0(t *testing.T) {
	c1 := (&config.InitContext{}).
		FromFile("conf-test-files/isa/20240901-1.yaml").Load()
	c1.PrintJson("ORIGINAL")

	c1.TheIsa()
	c1.PrintJson("RESULT")
}

func Test_Isa_2_0(t *testing.T) {
	c1 := (&config.InitContext{}).
		FromFile("conf-test-files/isa/20240901-2.yaml").Load()
	c1.PrintJson("ORIGINAL")

	c1.TheIsa()
	c1.PrintJson("RESULT")
}

func Test_Isa_3_0(t *testing.T) {
	c1 := (&config.InitContext{}).
		FromFile("conf-test-files/isa/20240901-3.yaml").Load()
	c1.PrintJson("ORIGINAL")

	c1.TheIsa()
	c1.PrintJson("RESULT")
}

func Test_Isa_4_0(t *testing.T) {
	c1 := (&config.InitContext{}).
		FromFile("conf-test-files/isa/20240901-4.yaml").Load()
	c1.PrintJson("ORIGINAL")

	c1.TheIsa()
	c1.PrintJson("RESULT")
}
