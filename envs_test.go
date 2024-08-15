package config_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_TranslateEnvs_KeySuffix(t *testing.T) {
	fmt.Println(config.TranslateEnvs_KeySuffix("Somethingp__superd__duperu__important"))
}
