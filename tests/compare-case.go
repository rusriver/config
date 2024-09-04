package main

import (
	"fmt"
	"testing"
)

type CompareCase struct {
	N      string
	a1, a2 any
}

func (cas *CompareCase) Do(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(cas.N, x)
		}
	}()
	if cas.a1 == cas.a2 {
		fmt.Println("OK  ", cas.N, ":", cas.a1, cas.a2)
	} else {
		fmt.Println("FAIL", cas.N, ":", cas.a1, cas.a2)
		t.Fail()
	}
}
