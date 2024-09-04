package config

import (
	"fmt"
	"testing"
)

func Test_Isa_1(t *testing.T) {
	type Case struct {
		N string
		p string
	}
	cases := []Case{
		{"1.0", "aa.ss"},
		{"1.5", ".aa.ss"},
		{"2.0", "^aa.ss"},
		{"3.0", "^.aa.ss"},
		{"4.0", "^^.aa.ss"},
		{"5.0", "^^^^.aa.ss"},
	}
	for _, cas := range cases {
		func() {
			defer func() {
				if x := recover(); x != nil {
					fmt.Println(cas.N, x)
				}
			}()
			n, p := (&TreeTraversalContext{}).howManyLevelsBackUp(cas.p)
			fmt.Println(cas.N, n, p)
		}()
	}
	/*
		1.0 0 aa.ss
		1.5 0 .aa.ss
		2.0 4114a72ea601 Invalid path syntax in $isa, '^aa.ss'
		3.0 1 aa.ss
		4.0 2 aa.ss
		5.0 4 aa.ss
	*/
}

func Test_Isa_2(t *testing.T) {
	a := []int{1, 2, 3}
	fmt.Println(a[0:0])
	a = []int{}
	fmt.Println(a[0:0])

	/*
	   []
	   []
	*/
}
