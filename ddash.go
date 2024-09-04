package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (c *Config) DDash() {
	ctx := &TreeTraversalContext{
		root:            c,
		antiLoopMap:     make(map[string]bool),
		currentLocation: &Location{},
	}
	c.DataSubTree = ctx.ddash_TreeTraversal("", c.DataSubTree)
}

func (ctx *TreeTraversalContext) ddash_TreeTraversal(nodeName string, node any) any {
	if len(nodeName) > 0 {
		ctx.currentLocation.Push(nodeName)
		defer func() {
			ctx.currentLocation.Pop()
		}()
	}
	// fmt.Println("++CURRLOC  ", *ctx.currentLocation)

	switch nv := node.(type) {
	case map[string]any:

		// traverse the tree depth-first
		for k, v := range nv {
			if k == "$-" {
				continue
			}
			nv[k] = ctx.ddash_TreeTraversal(k, v)
		}

		// back again - handle the $-
		if ddashObject, ok := nv["$-"]; ok {
			// fmt.Println("++DDASH  ", *ctx.currentLocation)

			switch ddashObject := ddashObject.(type) {
			case string:
				ctx.ddash_CmdExec(ddashObject)

			case []any:
				for _, ddashObject := range ddashObject {
					switch rawCmd := ddashObject.(type) {
					case string:
						ctx.ddash_CmdExec(rawCmd)
					}
				}

			} // switch

			delete(nv, "$-")
		}

	case []any:
		for n, e := range nv {
			nv[n] = ctx.ddash_TreeTraversal(strconv.Itoa(n), e)
		}

	} // switch

	return node
}

var reDdash01 = regexp.MustCompile(`\s+`)

func (ctx *TreeTraversalContext) ddash_CmdExec(raw string) {
	ctx.currentLocation.Push("$-")
	defer func() {
		ctx.currentLocation.Pop()
	}()
	// fmt.Println("++DDASH-CMD  ", *ctx.currentLocation)

	cmdArray := reDdash01.Split(raw, -1)
	switch cmdArray[0] {

	case "set":
		p := ctx.getAbsPath(cmdArray[1])
		ctx.root.Set(strings.Split(p, "."), cmdArray[2])

	case "string-interpolate":
		cmdArray = cmdArray[1:]
		for _, a := range cmdArray {
			p := ctx.getAbsPath(a)
			v := os.Expand(ctx.root.DotP(p).String(), func(s string) string {
				return ctx.root.DotP(s).String()
			})
			ctx.root.Set(strings.Split(p, "."), v)
		}

	case "append":
		p := ctx.getAbsPath(cmdArray[1])
		c := ctx.root.DotP(p)
		switch a := c.DataSubTree.(type) {
		case []any:
			a = append(a, cmdArray[2])
			ctx.root.Set(strings.Split(p, "."), a)
		}

	} // switch
}
