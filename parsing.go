package main

import (
	"image/color"
	"io"
	"strconv"
	"strings"
	"unicode"

	xq "github.com/antchfx/xquery/xml"
	"github.com/go-gl/mathgl/mgl64"
)

func makeXMLDoc(r io.Reader) *xq.Node {
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

func parseSvgString(s string) *svg {
	return parseSvgReader(strings.NewReader(s))
}

func parseSvgReader(r io.Reader) *svg {
	n := makeXMLDoc(r)
	return parseSvg(n)
}

func parseSvg(node *xq.Node) *svg {
	if node == nil {
		return nil
	}
	svg := new(svg)
	svg.elements = []svgElement{}

	for child := firstElementChild(node); child != nil; child = firstElementSibling(child) {
		if child.Data == "svg" {
			svg.width = parseFloat(child, "width")
			svg.height = parseFloat(child, "height")
			svg.elements = parseChildren(child)
		}
	}
	return svg
}

func parseChildren(node *xq.Node) []svgElement {
	children := []svgElement{}
	for child := firstElementChild(node); child != nil; child = firstElementSibling(child) {
		if child.Data == "line" {
			line := parseLine(child)
			line.data = parseElementData(child)
			children = append(children, &line)
		} else if child.Data == "rect" {
			r := parseRect(child)
			r.data = parseElementData(child)
			if r.dimension.X() == 0 && r.dimension.Y() == 0 {
				var p svgPoint
				p.data = r.data
				p.position = r.position
				children = append(children, &p)
			} else {
				children = append(children, &r)
			}
		} else if child.Data == "polygon" {
			polygon := parsePolygon(child)
			polygon.data = parseElementData(child)
			children = append(children, &polygon)
		} else if child.Data == "polyline" {
			polyline := parsePolyline(child)
			polyline.data = parseElementData(child)
			children = append(children, &polyline)
		} else if child.Data == "g" {
			group := parseGroup(child)
			group.data = parseElementData(child)
			children = append(children, &group)
		}
	}
	return children
}

func parseElementData(node *xq.Node) (data svgElementData) {
	if fill := node.SelectAttr("fill"); len(fill) > 0 {
		data.style.fillColor = parseColor(fill)
	}
	if fillOpacity := node.SelectAttr("fill-opacity"); len(fillOpacity) > 0 {
		fillOpacityF, _ := strconv.ParseFloat(fillOpacity, 32)
		data.style.fillColor.A = uint8(fillOpacityF * 255.0)
	}

	if stroke := node.SelectAttr("stroke"); len(stroke) > 0 {
		data.style.strokeColor = parseColor(stroke)
		strokeOpacity := node.SelectAttr("stroke-opacity")
		if len(strokeOpacity) > 0 {
			strokeOpacityF, _ := strconv.ParseFloat(strokeOpacity, 32)
			data.style.strokeColor.A = uint8(strokeOpacityF * 255.0)
		}
	}

	data.style.strokeWidth = parseFloat(node, "stroke-width")
	data.style.miterLimit = parseFloat(node, "stroke-miterlimit")

	// TODO: parse transform
	return
}

func parseGroup(node *xq.Node) (g svgGroup) {
	g.elements = parseChildren(node)
	return
}

func parseLine(node *xq.Node) (l svgLine) {
	l.from = mgl64.Vec2{parseFloat(node, "x1"), parseFloat(node, "y1")}
	l.to = mgl64.Vec2{parseFloat(node, "x2"), parseFloat(node, "y2")}
	return
}

func parseRect(node *xq.Node) (r svgRect) {
	r.position = mgl64.Vec2{parseFloat(node, "x"), parseFloat(node, "y")}
	r.dimension = mgl64.Vec2{parseFloat(node, "width"), parseFloat(node, "height")}
	return
}

func parsePolygon(node *xq.Node) (pg svgPolygon) {
	p := node.SelectAttr("points")
	pg.points = parsePoints(p)
	return
}

func parsePolyline(node *xq.Node) (pg svgPolyline) {
	p := node.SelectAttr("points")
	pg.points = parsePoints(p)
	return
}

func parsePoints(s string) []mgl64.Vec2 {
	points := []mgl64.Vec2{}
	splitted := strings.Split(s, " ")
	for _, pair := range splitted {
		i := strings.Index(pair, ",")
		if i != -1 {
			p1 := pair[:i]
			p2 := pair[i+1:]
			v := mgl64.Vec2{parseFloatString(p1), parseFloatString(p2)}
			points = append(points, v)
		}
	}
	return points
}

func goodRune(r rune) bool {
	return !(unicode.IsDigit(r) || r == '.')
}

func parseFloatString(s string) (v float64) {
	s = strings.TrimFunc(s, goodRune)

	if strings.Contains(s, ".") {
		v, _ = strconv.ParseFloat(s, 64)
	} else {
		i, _ := strconv.Atoi(s)
		v = float64(i)
	}
	return
}
func parseFloat(node *xq.Node, s string) (v float64) {
	s = node.SelectAttr(s)
	v = parseFloatString(s)
	return
}

func parseColor(s string) color.RGBA {
	if s == "none" {
		return color.RGBA{}
	}

	if s[0] == '#' {
		s = s[1:]
	}

	colInt, _ := strconv.ParseInt(s, 16, 32)
	var color color.RGBA
	color.R = uint8((colInt & 0xFF0000) >> 16)
	color.G = uint8((colInt & 0x00FF00) >> 8)
	color.B = uint8((colInt & 0x0000FF) >> 0)
	color.A = 0xFF
	return color
}
