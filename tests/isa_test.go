package main

import (
	"fmt"
	"testing"

	"github.com/rusriver/config/v2"
	"github.com/rusriver/nutz/controlflow"
	"github.com/stretchr/testify/require"
)

func Test_Isa_1_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/isa/20240901-1.yaml",
		"vbhs-_LYSsdYDtCtLwZ22gcr4_jLa_EBPz96VtwdIYY=",
	).Load()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")
	c1.TheIsa()
	c1.PrintJson("RESULT")

	// validation
	cases := []CompareCase{
		{"1.0", c1.DotP("objects.object-03.x").Int(), 1},
		{"2.0", c1.DotP("objects.object-01.a").Int(), 2},
		{"3.0", c1.DotP("objects.object-01.b").Int(), 3},
		{"4.0", c1.DotP("objects.object-01.c").Int(), 4},
	}
	for _, cas := range cases {
		cas.Do(t)
	}
}

func Test_Isa_2_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/isa/20240901-2.yaml",
		"wNb6xZmKUNTDfdgryp9zOeFfShKgKKQUU-_zcpCgvL8=",
	).Load()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")
	c1.TheIsa()
	c1.PrintJson("RESULT")

	// validation
	cases := []CompareCase{
		{"1.1", c1.DotP("objects.obj-01.a").Int(), 1},
		{"1.2", c1.DotP("objects.obj-01.b").Int(), 2},
		{"1.3", c1.DotP("objects.obj-01.name").String(), "obj-01"},
		{"2.1", c1.DotP("objects.obj-02.a").Int(), 1},
		{"2.2", c1.DotP("objects.obj-02.b").Int(), 2},
		{"2.3", c1.DotP("objects.obj-02.name").String(), "obj-02"},
		{"3.1", c1.DotP("objects.obj-03.a").Int(), 1},
		{"3.2", c1.DotP("objects.obj-03.b").Int(), 2},
		{"3.3", c1.DotP("objects.obj-03.name").String(), "obj-03"},
	}
	for _, cas := range cases {
		cas.Do(t)
	}
}

func Test_Isa_3_0(t *testing.T) {
	for i := 0; i < 50; i++ {
		c1 := (&config.InitContext{}).FromFile(
			"conf-test-files/isa/20240901-3.yaml",
			"BlElOym7C1y7fSarkeMJKBnzBnVY1k_yRtUE20tAnqk=",
		).Load()
		if i == 0 {
			fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
			c1.PrintJson("ORIGINAL")
		}

		c1.TheIsa()

		if i == 0 {
			c1.PrintJson("RESULT")
		}

		// validation - repeat several times, to make sure it works
		// on different map orderings
		cases := []CompareCase{
			{"1.0", c1.DotP("objects.oMultiple.a").Int(), 1},
			{"2.0", c1.DotP("objects.oMultiple.b").Int(), 20},
			{"3.0", c1.DotP("objects.oMultiple.c").Int(), 30},
			{"4.0", c1.DotP("objects.oMultiple.name").String(), "object-02"},
			{"5.0", c1.DotP("objects.oMultiple.o1_v").Int(), 15},
		}
		for _, cas := range cases {
			cas.Do(t)
		}
		fmt.Println("---")
		if t.Failed() {
			break
		}
	}
}

func Test_Isa_4_0(t *testing.T) {
	for i := 0; i < 50; i++ {
		c1 := (&config.InitContext{}).FromFile(
			"conf-test-files/isa/20240901-4.yaml",
			"zhP2hYQrqGiKBbfj0OH9WmxCniZuAHYOnhcNfP8BBio=",
		).Load()
		if i == 0 {
			fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
			c1.PrintJson("ORIGINAL")
		}

		c1.TheIsa()
		if i == 0 {
			c1.PrintJson("RESULT")
		}

		// validation - repeat several times, to make sure it works
		// on different map orderings
		cases := []CompareCase{
			{"1.1", c1.DotP("objects.array.0.a").Int(), 1},
			{"1.2", c1.DotP("objects.array.0.b").Int(), 2},
			{"1.3", c1.DotP("objects.array.0.name").String(), "object-01"},
			{"2.1", c1.DotP("objects.array.1.a").Int(), 1},
			{"2.2", c1.DotP("objects.array.1.b").Int(), 20},
			{"2.3", c1.DotP("objects.array.1.c").Int(), 30},
			{"2.4", c1.DotP("objects.array.1.name").String(), "object-02"},
		}
		for _, cas := range cases {
			cas.Do(t)
		}
		fmt.Println("---")
		if t.Failed() {
			break
		}
	}
}

func Test_Isa_5_0(t *testing.T) {
	c1 := (&config.InitContext{}).FromFile(
		"conf-test-files/isa/20250207-5.yaml",
		"E4IbXdtpjGTVga8YskqYVOgN_1jJ-1Si0BXG44_1XcU=",
	).Load().U()
	fmt.Println("CONFIG HASHES:", c1.InitContext.SourceHashesActual)
	c1.PrintJson("ORIGINAL")

	controlflow.Try(func() (err error) {
		c1.TheIsa()
		return
	}).Catch(func(e *controlflow.Exception) {
		require.Equal(t, " 4c82b5eaf5b88f268deb error: TheIsa() called without any kind of error handling", e.String())
	})

	controlflow.Try(func() (err error) {
		c1.UnU().TheIsa()
		return
	}).Catch(func(e *controlflow.Exception) {
		require.Equal(t, " e6269413b51b nonexistent map key 'objects.default'", e.String())
	})

	c1.PrintJson("RESULT")
}
