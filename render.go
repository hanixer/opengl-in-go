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
	return edge.a*x+edge.b*y+edge.c < 0
}

func (edge edgeEquation) evaluate(x, y float64) float64 {
	return edge.a*x + edge.b*y + edge.c
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

func makeColor(r, g, b float64) color.Color {
	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

func drawTriangle(v0, v1, v2 vertex, img draw.Image, c color.Color) {
	minX, minY, maxX, maxY := boundingBox(v0, v1, v2)

	e0 := newEdgeEquation(v1, v2)
	e1 := newEdgeEquation(v2, v0)
	e2 := newEdgeEquation(v0, v1)

	k0 := 1.0 / e0.evaluate(v0.x, v0.y)
	k1 := 1.0 / e1.evaluate(v1.x, v1.y)
	k2 := 1.0 / e2.evaluate(v2.x, v2.y)

	for y := minY + 0.5; y <= maxY+0.5; y++ {
		for x := minX + 0.5; x <= maxX+0.5; x++ {
			// compute baricentric coordinates
			w0 := k0 * e0.evaluate(x, y)
			w1 := k1 * e1.evaluate(x, y)
			w2 := k2 * e2.evaluate(x, y)

			if w0 > 0 && w1 > 0 && w2 > 0 {
				r := v0.r*w0 + v1.r*w1 + v2.r*w2
				g := v0.g*w0 + v1.g*w1 + v2.g*w2
				b := v0.b*w0 + v1.b*w1 + v2.b*w2
				colo := makeColor(r, g, b)
				img.Set(int(x), int(y), colo)
			}
		}
	}
}

func drawLine(x0, y0, x1, y1 int, img draw.Image) {
	y := y0
	edge := newEdgeEquation(vertex{x: float64(x0), y: float64(y0)}, vertex{x: float64(x1), y: float64(y1)})

	for x := x0; x <= x1; x++ {
		img.Set(x, y, color.Black)
		if !edge.test(float64(x), float64(y)) {
			y++
		}
	}
}
