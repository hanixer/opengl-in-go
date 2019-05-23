package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/color"
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
	strokeColor color.RGBA
	fillColor   color.RGBA
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

func (g *svgRect) isSvg() {}

type svgPolygon struct {
	data   svgElementData
	points []mgl64.Vec2
}

func (g *svgPolygon) isSvg() {}

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
