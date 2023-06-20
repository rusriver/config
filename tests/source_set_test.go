package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/rusriver/config/v2"
)

func Test_20230620_5(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		E(&err).
		Load()

	asd := conf.P("a", "s", "d")
	asd.Set([]string{"qq", "ww", "3", "congrads"}, "HELLO THERE")

	s := config.NewSource(func(opts *config.NewSource_Options) {
		opts.Config = conf
		opts.CommandBufferSize = 10
		opts.UpdatePeriod = time.Second * 1
	})

	s.Config.PrintJson("INIT")

	// TAKE NOTE: we must re-take the asd, because the former one hasn't
	// the Source set
	asd = conf.P("a", "s", "d")

	asd.P("N1").Set(nil, 555)
	s.Config.PrintJson("immediately after first modification")
	time.Sleep(time.Second * 2)
	s.Config.PrintJson("after delay")

	asd.P("N1").Set(nil, 111)
	asd.P("N2").Set(nil, 222)
	asd.P("N3").Set(nil, 333)
	asd.P("N4").Set(nil, 444)
	s.Config.PrintJson("immediately after batch")
	time.Sleep(time.Second * 4)
	s.Config.PrintJson("after delay")

	fmt.Println("MAKE A BIGGER BATCH")

	for x := 0; x < 30; x++ {
		asd.P("ARRAY", strconv.Itoa(x)).Set(nil, x)
	}
	s.Config.PrintJson("immediately after batch")
	time.Sleep(time.Second * 4)
	s.Config.PrintJson("after delay")
}

func Test_20230620_5_2(t *testing.T) {
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/c2.yaml").
		E(&err).
		Load()

	asd := conf.P("a", "s", "d")
	asd.Set([]string{"qq", "ww", "3", "congrads"}, "HELLO THERE")

	s := config.NewSource(func(opts *config.NewSource_Options) {
		opts.Config = conf
		opts.CommandBufferSize = 10
		opts.UpdatePeriod = time.Second * 1
	})

	s.Config.PrintJson("INIT")

	// TAKE NOTE: we must re-take the asd, because the former one hasn't
	// the Source set
	asd = conf.P("a", "s", "d")

	go func() {
		for x := 0; x < 300; x++ {
			asd.P("ARRAY-1", strconv.Itoa(x)).Set(nil, x)
		}
	}()

	go func() {
		for x := 0; x < 300; x++ {
			asd.P("ARRAY-2", strconv.Itoa(x)).Set(nil, x)
		}
	}()

	go func() {
		for x := 0; x < 300; x++ {
			asd.P("ARRAY-3", strconv.Itoa(x)).Set(nil, x)
		}
	}()

	time.Sleep(time.Second * 4)
	s.Config.PrintJson("after delay")
}

func Test_20230620_6(t *testing.T) {
	a := make([]int, 6)
	var aPtr *[]int = &a
	fmt.Printf("+++1 %+v\n", *aPtr)

	*aPtr = a
	fmt.Printf("+++2 %+v\n", *aPtr)

	a[3] = 3
	fmt.Printf("+++3 %+v\n", *aPtr)

	a = append(a, make([]int, 3)...)
	*aPtr = a
	fmt.Printf("+++4 %+v\n", *aPtr)

	a[3] = 33
	fmt.Printf("+++5 %+v\n", *aPtr)

	a[8] = 8
	fmt.Printf("+++6 %+v\n", *aPtr)
}
func Test_20230620_6_2(t *testing.T) {
	m := make(map[string]int)
	m["k1"] = 10

	m2 := m

	m["k2"] = 20

	fmt.Printf("+++1 %+v\n", m)
	fmt.Printf("+++2 %+v\n", m2)
}
