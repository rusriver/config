package main

import (
	"fmt"
	"testing"

	"github.com/rusriver/config/v2"
	"github.com/rusriver/nutz/controlflow"
)

func Test_Raisable_1_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/ddash/20240904-1.yaml",
		"FstTSOp-9WeXmf0h-RFB9rXVPEcnM43_JD1lfBEuga8=",
	).Load().U()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)

	c1 = c1.RaisableWithTag("RWT0ee3c018dcea")
	c1.PrintJson("PRINTOUT")

	controlflow.Try(func() (err error) {
		fmt.Println(c1.DotP("d.x1.b-non-existent").Int())
		return
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println(e)
	})
}
