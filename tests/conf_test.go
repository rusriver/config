package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/rusriver/config/v2"
	"github.com/rusriver/config/v2/deepcopy"
)

func Test_ConfigInheritance_1_0(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/config.yaml").
		Err(&err).
		LoadWithParenting()

	conf.PrintJson("1")
}

func Test_20230620_1(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		Err(&err).
		Load()

	conf.PrintJson("1")

	asd := conf.P("a", "s", "d")
	v1 := asd.P("f", "g", "h")
	fmt.Println(v1.Int())

	as := conf.P("a", "s")
	dfg := as.P("d", "f", "g")
	v2 := dfg.P("h")
	fmt.Println(v2.Int())

	fmt.Printf("%+v\n", v1.GetCurrentLocationPlusPath())
	fmt.Printf("%+v\n", v2.GetCurrentLocationPlusPath())
	fmt.Printf("%+v\n", v2.GetCurrentLocationPlusPath("z", "x", "c"))
}

func Test_20230620_2(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		Err(&err).
		Load()

	conf.PrintJson("initial")

	asd := conf.P("a", "s", "d")
	asd.Set([]string{"qq", "ww", "3", "congrads"}, "HELLO THERE")

	conf.PrintJson("after modification")
}

func Test_20230620_3(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		Err(&err).
		Load()

	conf.PrintJson("initial")

	asd := conf.P("a", "s", "d")
	asd.Set([]string{"qq", "ww", "3", "congrads"}, "HELLO THERE")

	conf.PrintJson("after modification")

	conf2 := deepcopy.Copy(conf).(*config.Config) // won't work with Source though

	conf2.PrintJson("after deepcopy")
}

func Test_20230620_4(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		Err(&err).
		Load()

	conf.PrintJson("initial")

	asd := conf.P("a", "s", "d")

	asd.Err(&err).Set([]string{"ARRAY", "3"}, 3)
	if err != nil {
		t.Fatalf("%v", err)
	}
	asd.Err(&err).Set([]string{"ARRAY", "6"}, 6)
	if err != nil {
		t.Fatalf("%v", err)
	}

	for x := 5; x < 30; x++ {
		asd.Err(&err).Set([]string{"ARRAY", strconv.Itoa(x + 3)}, x-5)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}

	conf.PrintJson("conf after modification")
	asd.PrintJson("asd after modification")
}
