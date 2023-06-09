package main

import (
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_ConfigInheritance_1_0(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/config.yaml").
		E(&err).
		LoadWithParenting()

	conf.PrintJson("1")
}
