package main

import (
	"fmt"
	"os"

	xq "github.com/antchfx/xquery/xml"
)

// Scene is a scene. It contains other objects, meshes, integrators, etc.
type scene struct {
}

// Object is a common superclass of all objects
type object interface {
	addChild(object)
	setParent(object)
	activate()
}

type property interface{}

type propertyList map[string]property

func (plist propertyList) setString(name, value string) {
	plist[name] = value
}

func createInstance(name string, plist propertyList) object {
	return &ttt{name: name, plist: plist}
}

func mappingfu(m propertyList) {
	fmt.Println(m["key"])
	m["key"] = 2
}

func main() {
	m := make(propertyList)
	m["key"] = 1
	mappingfu(m)
	fmt.Print(m)
	f, _ := os.Open("test.xml")
	n, _ := xq.Parse(f)
	parseIt(n)
}
