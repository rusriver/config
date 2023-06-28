package config

import (
	"encoding/json"
	"fmt"
)

func (c *Config) PrintJson(tag string) {
	bb, err := json.MarshalIndent(c, "", "    ")
	fmt.Printf("PRINT JSON, tag='%v' err='%v'        -- %v\n%s\n", tag, err, tag, bb)
}
