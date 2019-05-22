package main

import (
	"fmt"
	xq "github.com/antchfx/xquery/xml"
	"io"
)

func makeXmlDoc(r io.Reader) *xq.Node {
	n, _ := xq.Parse(r)
	return n
}

func firstElementChild(n *xq.Node) *xq.Node {
	n = n.FirstChild
	for {
		if n == nil {
			return nil
		}
		if n.Type == xq.ElementNode {
			return n
		}
		n = n.NextSibling
	}
}

func firstElementSibling(n *xq.Node) *xq.Node {
	for {
		n = n.NextSibling
		if n == nil {
			return nil
		}
		if n.Type == xq.ElementNode {
			return n
		}
	}
}

func parseSvg(node *xq.Node) {
	if node == nil {
		return
	}
	svg := new(svg)
	svg.elements = []svgElement{}
	elem := firstElementChild(node)
	for elem != nil {
		if elem.Data == "line" {

		}
	}
	fmt.Println(node.Data, node.Type)
	parseSvg(firstElementChild(node))
	parseSvg(firstElementSibling(node))
}

func parseElementData(node *xq.Node) (data svgElementData) {
	return
}
