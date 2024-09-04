package config

import (
	"fmt"
	"strconv"
	"strings"

	// "github.com/lithammer/shortuuid/v3"
	"github.com/rusriver/config/v2/deepcopy"
)

type theIsaContext struct {
	root            *Config
	antiLoopMap     map[string]bool
	currentLocation *Location
}

func (c *Config) TheIsa() {
	ctx := &theIsaContext{
		root:            c,
		antiLoopMap:     make(map[string]bool),
		currentLocation: &Location{},
	}
	ctx.resolveRelativePaths_TreeTraversal("", c.DataSubTree)
	// c.PrintJson("============== $isa middle") //--==
	c.DataSubTree = ctx.applyTheIsa_TreeTraversal("", c.DataSubTree)
}

func (ctx *theIsaContext) resolveRelativePaths_TreeTraversal(nodeName string, node any) any {
	if len(nodeName) > 0 {
		ctx.currentLocation.Push(nodeName)
		defer func() {
			ctx.currentLocation.Pop()
		}()
	}

	switch nv := node.(type) {
	case map[string]any:

		// traverse the tree depth-first
		for k, v := range nv {
			nv[k] = ctx.resolveRelativePaths_TreeTraversal(k, v)
		}

		// back again - handle the $isa
		if isaObject, ok := nv["$isa"]; ok {

			switch isaObject2 := isaObject.(type) {
			case string:
				nv["$isa"] = ctx.getAbsPath(isaObject2)

			case []any:
				paths := make([]any, 0, 8)
				for _, isaPath3 := range isaObject2 {
					switch isaPath4 := isaPath3.(type) {
					case string:
						paths = append(paths, ctx.getAbsPath(isaPath4))
					}
				}
				nv["$isa"] = paths

			} // switch
		}

	case []any:
		for n, e := range nv {
			nv[n] = ctx.resolveRelativePaths_TreeTraversal(strconv.Itoa(n), e)
		}

	} // switch

	return node
}

func (ctx *theIsaContext) applyTheIsa_TreeTraversal(nodeName string, node any) any {
	if len(nodeName) > 0 {
		ctx.currentLocation.Push(nodeName)
		defer func() {
			ctx.currentLocation.Pop()
		}()
	}
	// fmt.Println("++  ", *ctx.currentLocation) //--==

	switch nv := node.(type) {

	case map[string]any:

		// traverse the tree depth-first
		for k, v := range nv {
			nv[k] = ctx.applyTheIsa_TreeTraversal(k, v)
		}

		// back again - handle the $isa
		if isaObject, ok := nv["$isa"]; ok {

			paths := make([]string, 0, 8)

			switch isaObject2 := isaObject.(type) {
			case string:
				paths = append(paths, isaObject2)

			case []any:
				for _, isaPath3 := range isaObject2 {
					switch isaPath4 := isaPath3.(type) {
					case string:
						paths = append(paths, isaPath4)
					}
				}

			}
			node = ctx.applyTheIsa_DoMultipleInheritance(paths, node)
			switch nv := node.(type) {
			case map[string]any:
				delete(nv, "$isa")
			}
			node = ctx.applyTheIsa_TreeTraversal(nodeName, node) // repeat itself until there's no $isa left
		}

	case []any:
		for n, e := range nv {
			nv[n] = ctx.applyTheIsa_TreeTraversal(strconv.Itoa(n), e)
		}

	} // switch

	return node
}

func (ctx *theIsaContext) applyTheIsa_DoMultipleInheritance(isaPaths []string, aLast any) any {
	p1 := isaPaths[0]
	isaPaths = isaPaths[1:]

	if ctx.antiLoopMap[p1] {
		panic(fmt.Errorf("d2e955d2-82, $isa loop; %v", strings.Join(*ctx.currentLocation, "->")))
	}
	ctx.antiLoopMap[p1] = true

	defer delete(ctx.antiLoopMap, p1)

	// get first object
	// uuid := shortuuid.New()
	newBaseObject := deepcopy.Copy(ctx.root.DotP(p1).DataSubTree)
	// fmt.Println("FIRST OBJECT", uuid, p1, *ctx.currentLocation)               //--==
	// (&Config{DataSubTree: newBaseObject}).PrintJson("FIRST OBJECT 1 " + uuid) //--==
	newBaseObject = ctx.applyTheIsa_TreeTraversal(p1, newBaseObject)
	// fmt.Println("FIRST OBJECT", uuid, p1, *ctx.currentLocation)               //--==
	// (&Config{DataSubTree: newBaseObject}).PrintJson("FIRST OBJECT 2 " + uuid) //--==

	for _, pN := range isaPaths { // all except the first one

		if ctx.antiLoopMap[pN] {
			panic(fmt.Errorf("d2e955d2-97, $isa loop; %v", strings.Join(*ctx.currentLocation, "->")))
		}
		ctx.antiLoopMap[pN] = true

		newBaseObject = Extend_v2_any(newBaseObject, ctx.root.DotP(pN).DataSubTree)
		newBaseObject = ctx.applyTheIsa_TreeTraversal(pN, newBaseObject)

		delete(ctx.antiLoopMap, pN)
	}

	newBaseObject = Extend_v2_any(newBaseObject, aLast)

	// fmt.Println("FIRST OBJECT", uuid, p1, *ctx.currentLocation)                 //--==
	// (&Config{DataSubTree: newBaseObject}).PrintJson("FIRST OBJECT 165 " + uuid) //--==

	return newBaseObject
}

func (ctx *theIsaContext) getAbsPath(p string) string {
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

func (ctx *theIsaContext) howManyLevelsBackUp(p string) (n int, p2 string) {
	p2 = p
	for _, c := range p {
		switch c {
		case '^':
			n++
		case '.':
			break
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
