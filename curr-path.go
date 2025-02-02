package config

import "strings"

type CurrPath []string

func (c *CurrPath) String() string {
	return strings.Join(*c, ".")
}

func (c *CurrPath) Push(p string) {
	*c = append(*c, p)
}

func (c *CurrPath) Pop() string {
	p := (*c)[len(*c)-1]
	*c = (*c)[:len(*c)-1]
	return p
}
