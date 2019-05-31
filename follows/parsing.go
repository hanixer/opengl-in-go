package main

import (
	"fmt"
	"log"

	xq "github.com/antchfx/xquery/xml"
)

const (
	sceneTag int = iota
	meshTag
	bsdfTag
	phaseFunction
	emitter
	medium
	camera
	integrator
	sampler
	testTag
	reconstructionFilter
	booleanTag
	integerTag
	floatTag
	stringTag
	pointTag
	vectorTag
	colorTag
	transformTag
	translateTag
	matrixTag
	rotateTag
	scaleTag
	lookatTag
)

var tags = map[string]int{
	"scene":      sceneTag,
	"mesh":       meshTag,
	"bsdf":       bsdfTag,
	"phase":      phaseFunction,
	"emitter":    emitter,
	"medium":     medium,
	"camera":     camera,
	"integrator": integrator,
	"sampler":    sampler,
	"test":       testTag,
	"rfilter":    reconstructionFilter,
	"boolean":    booleanTag,
	"integer":    integerTag,
	"float":      floatTag,
	"string":     stringTag,
	"point":      pointTag,
	"vector":     vectorTag,
	"color":      colorTag,
	"transform":  transformTag,
	"translate":  translateTag,
	"matrix":     matrixTag,
	"rotate":     rotateTag,
	"scale":      scaleTag,
	"lookat":     lookatTag,
}

func parseIt(node *xq.Node) {
	for node = node.FirstChild; node != nil; node = node.NextSibling {
		if node.Data != "xml" {
			break
		}
	}
	parseTag(node, make(propertyList))
}

func isObject(tag int) bool {
	switch tag {
	case sceneTag, meshTag, bsdfTag, phaseFunction, emitter, medium, camera:
		return true
	default:
		return false
	}
}

type ttt struct {
	name  string
	plist propertyList
}

func (t ttt) activate() {
	fmt.Println("activate! this -", t.name, "props", t.plist)
}

func (t ttt) addChild(o object) {
	fmt.Println("add child", o, ". this -", t.name)
}

func (t ttt) setParent(parent object) {}

func parseTag(node *xq.Node, plist propertyList) (result object) {

	children := []object{}

	currPList := make(propertyList)

	tag, ok := tags[node.Data]
	if !ok {
		log.Fatalln("tag", node.Data, "not found", tag)
	}

	for chNode := node.FirstChild; chNode != nil; chNode = chNode.NextSibling {
		if chNode.Type == xq.ElementNode {
			child := parseTag(chNode, currPList)
			if child != nil {
				children = append(children, child)
			}
		}
	}

	if isObject(tag) {
		result = createInstance(node.SelectAttr("type"), currPList)

		for _, child := range children {
			result.addChild(child)
			child.setParent(result)
		}

		result.activate()
	} else {
		parseProperty(node, tag, plist)
	}

	return nil
}

func parseProperty(node *xq.Node, tag int, plist propertyList) {
	switch tag {
	case stringTag:
		plist[node.SelectAttr("name")] = node.SelectAttr("value")
	case floatTag:
		strconv.ParseFloat(node.SelectAttr("value"))
		plist[node.SelectAttr("name")] = node.SelectAttr("value")
	}
}
