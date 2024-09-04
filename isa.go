package config

import (
	"fmt"
	"strconv"
	"strings"

	// "github.com/lithammer/shortuuid/v3"
	"github.com/rusriver/config/v2/deepcopy"
)

func (c *Config) TheIsa() {
	ctx := &TreeTraversalContext{
		root:            c,
		antiLoopMap:     make(map[string]bool),
		currentLocation: &Location{},
	}
	ctx.resolveRelativePaths_TreeTraversal("", c.DataSubTree)
	// c.PrintJson("============== $isa middle") //--==
	c.DataSubTree = ctx.applyTheIsa_TreeTraversal("", c.DataSubTree)
}

func (ctx *TreeTraversalContext) resolveRelativePaths_TreeTraversal(nodeName string, node any) any {
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
			if k == "$isa" {
				continue
			}
			nv[k] = ctx.resolveRelativePaths_TreeTraversal(k, v)
		}

		// back again - handle the $isa
		if isaObject, ok := nv["$isa"]; ok {

			switch isaObject := isaObject.(type) {
			case string:
				nv["$isa"] = ctx.getAbsPath(isaObject)

			case []any:
				paths := make([]any, 0, 8)
				for _, isaPath := range isaObject {
					switch isaPath := isaPath.(type) {
					case string:
						paths = append(paths, ctx.getAbsPath(isaPath))
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

func (ctx *TreeTraversalContext) applyTheIsa_TreeTraversal(nodeName string, node any) any {
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
			if k == "$isa" {
				continue
			}
			nv[k] = ctx.applyTheIsa_TreeTraversal(k, v)
		}

		// back again - handle the $isa
		if isaObject, ok := nv["$isa"]; ok {

			paths := make([]string, 0, 8)

			switch isaObject := isaObject.(type) {
			case string:
				paths = append(paths, isaObject)

			case []any:
				for _, isaPath := range isaObject {
					switch isaPath := isaPath.(type) {
					case string:
						paths = append(paths, isaPath)
					}
				}

			}
			node = ctx.applyTheIsa_DoMultipleInheritance(paths, node)
			switch node := node.(type) {
			case map[string]any:
				delete(node, "$isa")
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

func (ctx *TreeTraversalContext) applyTheIsa_DoMultipleInheritance(isaPaths []string, aLast any) any {
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
