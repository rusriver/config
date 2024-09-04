package main

import (
	"fmt"
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_DDash_1_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/ddash/20240904-1.yaml",
		"FstTSOp-9WeXmf0h-RFB9rXVPEcnM43_JD1lfBEuga8=",
	).Load().U()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")
	c1.DDash()
	c1.PrintJson("RESULT")

	// validation
	cases := []CompareCase{
		{"1.0", c1.DotP("o.array.4.vv1").Int(), 77},
		{"2.0", c1.DotP("o.array.4.vv2").Int(), 88},
		{"3.0", c1.DotP("o.array.4.name").String(), "new-object-03"},
	}
	for _, cas := range cases {
		cas.Do(t)
	}
}

func Test_DDash_2_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/ddash/20240904-2.yaml",
		"BtZSmKsmsyYUDImOR9hA_v_cxbqaVF9dUavD0VtS530=",
	).Load().U()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")
	c1.DDash()
	c1.PrintJson("RESULT")

	// validation
	cases := []CompareCase{
		{"1.0", c1.DotP("bb.x1.$isa").String(), "z.x.c1.BUTTERFLY"},
		{"2.0", c1.DotP("bb.x2.$isa").String(), "z.x.c2.BUTTERFLY"},
	}
	for _, cas := range cases {
		cas.Do(t)
	}
}

func Test_DDash_3_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/ddash/20240904-3.yaml",
		"U68CqFT2nL8gJvVRxRPlgMi2IH2n4rD4CemTyNYRSTY=",
	).Load().U()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")
	c1.DDash()
	c1.PrintJson("RESULT")

	// validation
	cases := []CompareCase{
		{"1.0", c1.DotP("array.0").String(), "host1.com"},
		{"2.0", c1.DotP("array.1").String(), "host2.com"},
		{"3.0", c1.DotP("array.2").String(), "host3.com"},
		{"4.0", c1.DotP("array.3").String(), "host53.com"},
		{"5.0", c1.DotP("array.4").String(), "host54.com"},
	}
	for _, cas := range cases {
		cas.Do(t)
	}
}
