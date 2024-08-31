package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_TranslateEnvs_KeySuffix(t *testing.T) {
	fmt.Println(config.TranslateEnvs_KeySuffix("Somethingp__superd__duperu__important"))
}

func Test_Envs02(t *testing.T) {
	type Case struct {
		N     string
		Parts []string
	}
	cases := []Case{
		{"1.0", []string{"asd-xxx", "x15", "$isa"}},
	}
	for _, cas := range cases {
		k := strings.Join(cas.Parts, "_")
		k = strings.ToUpper(k)
		k = config.ReEnvs01.ReplaceAllString(k, "")
		fmt.Println(cas.N, k)
	}
}
