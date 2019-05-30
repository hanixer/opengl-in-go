package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

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

func (g *svgPolyline) isSvg() {}

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

func drawSvg(svg *svg, samplesCount int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(svg.width), int(svg.height)))
	now := time.Now()
	drawElements(svg.elements, img, samplesCount)
	fmt.Println("render time", time.Now().Sub(now).Seconds())
	return img
}

func drawElements(elements []svgElement, img draw.Image, samplesCount int) {
	for _, elem := range elements {
		switch v := elem.(type) {
		case *svgLine:
			drawLineF(v.from.X(), v.from.Y(), v.to.X(), v.to.Y(), img, v.data.style.strokeColor)
		case *svgPoint:
			drawPoint(v.position.X(), v.position.Y(), img, v.data.style.strokeColor)
		case *svgPolygon:
			drawSvgPolygon(v, img, samplesCount)
		case *svgPolyline:
			drawSvgPolyline(v, img, samplesCount)
		case *svgGroup:
			drawElements(v.elements, img, samplesCount)
		}
	}
}

func drawSvgPolyline(polyline *svgPolyline, img draw.Image, samplesCount int) {
	for i := 0; i < len(polyline.points); i++ {
		p0 := polyline.points[i]
		p1 := polyline.points[(i+1)%len(polyline.points)]
		drawLineF(p0.X(), p0.Y(), p1.X(), p1.Y(), img, polyline.data.style.strokeColor)
	}
}

func drawSvgPolygon(polygon *svgPolygon, img draw.Image, samplesCount int) {
	triangles := triangulate(polygon.points)
	for i := 0; i+2 < len(triangles); i += 3 {
		fillTriangle(triangles[i], triangles[i+1], triangles[i+2], img, polygon.data.style.fillColor, samplesCount)

		// drawLinePoints(triangles[i], triangles[i+1], img, polygon.data.style.strokeColor)
		// drawLinePoints(triangles[i+1], triangles[i+2], img, polygon.data.style.strokeColor)
	}

	if polygon.data.style.strokeColor.A != 0 {
		ps := polygon.points
		for i := 0; i < len(ps); i++ {
			// fmt.Println(ps[i], ps[(i+1)%len(ps)])
			// drawLinePoints(ps[i], ps[(i+1)%len(ps)], img, polygon.data.style.strokeColor)
		}
	}
}
