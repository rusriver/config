package config

import (
	"fmt"
	"strconv"
	"strings"

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
	ctx.theIsa_traverseTheTree("", c.DataSubTree)
}

func (ctx *theIsaContext) theIsa_traverseTheTree(nodeName string, node any) any {
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
			nv[k] = ctx.theIsa_traverseTheTree(k, v)
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
			node = ctx.theIsa_handleMultipleInheritance(paths, node)
			switch nv := node.(type) {
			case map[string]any:
				delete(nv, "$isa")
			}
			node = ctx.theIsa_traverseTheTree(nodeName, node) // repeat itself until there's no $isa left
		}

	case []any:
		for n, e := range nv {
			nv[n] = ctx.theIsa_traverseTheTree(strconv.Itoa(n), e)
		}

	} // switch

	return node
}

func (ctx *theIsaContext) theIsa_handleMultipleInheritance(isaPaths []string, aLast any) any {
	p1 := isaPaths[0]
	isaPaths = isaPaths[1:]

	if strings.HasPrefix(p1, "^.") {
		p1 = p1[2:]
		p1 = strings.Join((*ctx.currentLocation)[:len(*ctx.currentLocation)-1], ".") + "." + p1
	}
	if ctx.antiLoopMap[p1] {
		panic(fmt.Errorf("d2e955d2-82, $isa loop; %v", strings.Join(*ctx.currentLocation, "->")))
	}
	ctx.antiLoopMap[p1] = true

	defer delete(ctx.antiLoopMap, p1)

	newBaseObject := deepcopy.Copy(ctx.root.DotP(p1).DataSubTree)
	ctx2 := &theIsaContext{
		root:            ctx.root,
		antiLoopMap:     ctx.antiLoopMap,
		currentLocation: &Location{},
	}
	ctx2.theIsa_traverseTheTree("", newBaseObject)

	for _, pN := range isaPaths { // all except the first one

		if strings.HasPrefix(pN, "^.") {
			pN = pN[2:]
			pN = strings.Join((*ctx.currentLocation)[:len(*ctx.currentLocation)-1], ".") + "." + pN
		}
		if ctx.antiLoopMap[pN] {
			panic(fmt.Errorf("d2e955d2-97, $isa loop; %v", strings.Join(*ctx.currentLocation, "->")))
		}
		ctx.antiLoopMap[pN] = true

		newBaseObject = Extend_v2_any(newBaseObject, ctx.root.DotP(pN).DataSubTree)
		ctx2 := &theIsaContext{
			root:            ctx.root,
			antiLoopMap:     ctx.antiLoopMap,
			currentLocation: &Location{},
		}
		ctx2.theIsa_traverseTheTree("", newBaseObject)

		delete(ctx.antiLoopMap, pN)
	}

	newBaseObject = Extend_v2_any(newBaseObject, aLast)

	return newBaseObject
}
