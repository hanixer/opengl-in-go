package main

import (
	"image/color"
	"image/draw"
)

type vertex struct {
	x float64
	y float64
	r float64
	g float64
	b float64
}

func edge(x, y float64) {

}

type edgeEquation struct {
	a, b, c float64
}

// v1 - v0
// n2 = -(v1y - v0y), (v1x - v0x)
// p0 = v0
// n * ((x, y) - v0)
// -(v1y - v0y) * (x - v0x) + (v1x - v0x) * (y - v0y)
// -(v1y - v0y) * x - -(v1y - v0y) * v0x + (v1x - v0x) * y - (v1x - v0x) * v0y
// a = -(v1y - v0y)
// b = (v1x - v0x)
// c = - -(v1y - v0y) * v0x - (v1x - v0x) * v0y
func newEdgeEquation(v0, v1 vertex) edgeEquation {
	var eq edgeEquation
	eq.a = v0.y - v1.y
	eq.b = v1.x - v0.x
	eq.c = -(eq.a*v0.x + eq.b*v0.y)
	return eq
}

// v0 = (0, 0)
// v1 = (5, 10)
// n = (-10, 5)
// (-10, 5) * (5, 5) = -10 * 5 + 25 =
func (edge edgeEquation) test(x, y float64) bool {
	return edge.a*x+edge.b*y+edge.c > 0
}

type parameterEquation struct {
	a, b, c float64
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func boundingBox(v0, v1, v2 vertex) (x0, y0, x1, y1 float64) {
	x0 = min(min(v0.x, v1.x), v2.x)
	y0 = min(min(v0.y, v1.y), v2.y)
	x1 = max(max(v0.x, v1.x), v2.x)
	y1 = max(max(v0.y, v1.y), v2.y)
	return
}

func drawTriangle(v0, v1, v2 vertex, img draw.Image, c color.Color) {
	minX, minY, maxX, maxY := boundingBox(v0, v1, v2)

	e0 := newEdgeEquation(v1, v2)
	e1 := newEdgeEquation(v2, v0)
	e2 := newEdgeEquation(v0, v1)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if e0.test(x, y) && e1.test(x, y) && e2.test(x, y) {
				img.Set(int(x), int(y), c)
			}
		}
	}
}
