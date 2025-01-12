package config

import (
	"fmt"
	"strings"
)

type Location []string

func (loc *Location) Push(s string) {
	*loc = append(*loc, s)
}

func (loc *Location) Pop() {
	*loc = (*loc)[:len(*loc)-1]
}

type TreeTraversalContext struct {
	root            *Config
	antiLoopMap     map[string]bool
	currentLocation *Location
}

func (ctx *TreeTraversalContext) getAbsPath(p string) string {
	n, p := ctx.howManyLevelsBackUp(p)
	len_cl := len(*ctx.currentLocation)
	if n > 0 && len_cl > 0 {
		if n > len_cl {
			n = len_cl
		}
		p = strings.Join((*ctx.currentLocation)[:len_cl-n], ".") + "." + p
	}
	return p
}

func (ctx *TreeTraversalContext) howManyLevelsBackUp(p string) (n int, p2 string) {
	p2 = p
L:
	for _, c := range p {
		switch c {
		case '^':
			n++
		case '.':
			break L
		}
	}
	p = p[n:]
	if n > 0 && len(p) > 0 {
		if p[0] == '.' {
			p = p[1:]
		} else {
			panic(fmt.Errorf("4114a72ea601 Invalid path syntax in $isa, '%v'", p2))
		}
	}
	return n, p
}
