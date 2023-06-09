package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rusriver/config/v2"
)

func Test_ConfigInheritance_1_0(t *testing.T) {
	result := (&config.InitContext{}).FromFile("conf-test-files/config.yaml").LoadWithParenting()

	bb, _ := json.MarshalIndent(result.Root, "", "    ")
	fmt.Println(string(bb))
}
