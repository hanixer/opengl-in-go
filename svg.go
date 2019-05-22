package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/color"
	"image/draw"
)

// svgelementtype represents type of svg document element
type svgElementType int

// element types enum
const (
	none svgElementType = iota
	point
	line
	polyline
	rect
	polygon
	ellipse
	imageType
	group
)

// style represents the style of drawing
type style struct {
	strokeColor color.Color
	fillColor   color.Color
	strokeWidth float64
	miterLimit  float64
}

type svgElement interface {
	isSvg()
}

type svgElementData struct {
	svgType   svgElementType
	style     style
	transform mgl64.Mat3
}

type svgGroup struct {
	data     svgElementData
	elements []svgElement
}

func (g *svgGroup) isSvg() {}

type svgPoint struct {
	data     svgElementData
	position mgl64.Vec2
}

func (g *svgPoint) isSvg() {}

type svgLine struct {
	data svgElementData
	from mgl64.Vec2
	to   mgl64.Vec2
}

func (g *svgLine) isSvg() {}

type svgPolyline struct {
	data   svgElementData
	points []mgl64.Vec2
}

type svgRect struct {
	data      svgElementData
	position  mgl64.Vec2
	dimension mgl64.Vec2
}

type svgPolygon struct {
	data   svgElementData
	points []mgl64.Vec2
}

type svgEllipse struct {
	data   svgElementData
	center mgl64.Vec2
	radius mgl64.Vec2
}

type svgImage struct {
	data      svgElementData
	position  mgl64.Vec2
	dimension mgl64.Vec2
	buffer    image.Image
}

type svg struct {
	width, height float64
	elements      []svgElement
}

func drawSvg(target draw.Image, svgElem svgElement) {
	switch v := svgElem.(type) {
	case *svgLine:
		drawLineF(v.from.X(), v.from.Y(), v.to.X(), v.to.Y(), target)
	case *svgGroup:
		for _, el := range v.elements {
			drawSvg(target, el)
		}
	}
}

func ononon() {
	var jjj svgElement
	lin := svgLine{}
	jjj = &lin
	switch v := jjj.(type) {
	case *svgLine:
		fmt.Println(v.data)
	}
	fmt.Println(jjj)
}
