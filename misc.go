package config

import (
	"encoding/json"
	"fmt"
)

func (c *Config) PrintJson(tag string) {
	bb, err := json.MarshalIndent(c, "", "    ")
	fmt.Printf("tag='%v' err='%v'\n%s\n", tag, err, bb)
}
