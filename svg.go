package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/mathgl/mgl64"
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

func drawSvg(svg *svg) draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(svg.width), int(svg.height)))
	drawElements(svg.elements, img)

	return img
}

func drawElements(elements []svgElement, img draw.Image) {
	fmt.Println(len(elements))
	for _, elem := range elements {
		switch v := elem.(type) {
		case *svgLine:
			drawLineF(v.from.X(), v.from.Y(), v.to.X(), v.to.Y(), img, v.data.style.strokeColor)
		case *svgPoint:
			drawPoint(v.position.X(), v.position.Y(), img, v.data.style.strokeColor)
		case *svgPolygon:
			drawSvgPolygon(v, img)
		case *svgGroup:
			drawElements(v.elements, img)
		}
	}
}

func drawSvgPolygon(polygon *svgPolygon, img draw.Image) {
	triangles := triangulate(polygon.points)

	for i := 0; i+2 < len(triangles); i += 3 {
		// drawTriangle2(triangles[i], triangles[i+1], triangles[i+2], img, polygon.data.style.fillColor)

		// drawLinePoints(triangles[i], triangles[i+1], img, polygon.data.style.strokeColor)
		// drawLinePoints(triangles[i+1], triangles[i+2], img, polygon.data.style.strokeColor)
	}

	ps := polygon.points
	for i := 0; i < len(polygon.points); i += 2 {
		drawLinePoints(ps[i], ps[(i+1)%len(ps)], img, polygon.data.style.strokeColor)
	}
}

func triangulate(points []mgl64.Vec2) []mgl64.Vec2 {
	result := []mgl64.Vec2{}
	for i := 1; i+1 < len(points); i += 2 {
		result = append(result, points[0])
		result = append(result, points[i])
		result = append(result, points[i+1])
	}

	return result
}
