package main

import (
	"fmt"
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

func newEdgeEquation(v0, v1 vertex) edgeEquation {
	var eq edgeEquation
	eq.a = v0.y - v1.y
	eq.b = v1.x - v0.x
	eq.c = v0.x*v1.y - v1.x*v0.y
	return eq
}
func newEdgeEquationInt(x0, y0, x1, y1 int) edgeEquation {
	var e edgeEquation
	e.a = float64(y0 - y1)
	e.b = float64(x1 - x0)
	e.c = float64(x0*y1 - x1*y0)
	return e
}

func (edge edgeEquation) test(x, y float64) bool {
	return edge.a*x+edge.b*y+edge.c > 0
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

func drawLineF(x0, y0, x1, y1 float64, img draw.Image) {
	drawLine(int(x0), int(y0), int(x1), int(y1), img)
}

func drawLine(x0, y0, x1, y1 int, img draw.Image) {
	steep := false
	incr := 1
	if x1-x0 == 0 {
		incr = 1
		steep = true
	} else if slope := float64(y1-y0) / float64(x1-x0); slope > 1.0 {
		incr = 1
		steep = true
	} else if slope >= 0.0 {
		incr = 1
		steep = false
	} else if slope > -1.0 {
		incr = -1
		steep = false
	} else if slope <= 1.0 {
		incr = -1
		steep = true
	}

	var u0, u1, v0, _, edge = chooseEndPoints(steep, x0, x1, y0, y1)

	v := v0
	for u := u0; u <= u1; u++ {
		if steep {
			img.Set(v, u, color.White)
			fmt.Println(edge.evaluate(float64(v+incr), float64(u)+1.5), edge.a, edge.b, edge.c, incr, v)
			if incr > 0 && edge.evaluate(float64(v+incr), float64(u)+1.5) > 0 {
				v += incr
			} else if incr < 0 && edge.evaluate(float64(v+incr), float64(u)+1.5) < 0 {
				v += incr
			}
		} else {
			img.Set(u, v, color.White)
			fmt.Println(edge.evaluate(float64(u)+1.5, float64(v+incr)), edge.a, edge.b, edge.c, incr, v)
			if incr > 0 && edge.evaluate(float64(u)+1.5, float64(v+incr)) < 0 {
				v += incr
			} else if incr < 0 && edge.evaluate(float64(u)+1.5, float64(v+incr)) > 0 {
				v += incr
			}
		}
	}
}

func chooseEndPoints(steep bool, x0, x1, y0, y1 int) (int, int, int, int, edgeEquation) {
	if steep {
		if y0 < y1 {
			return y0, y1, x0, x1, newEdgeEquationInt(x0, y0, x1, y1)
		}
		return y1, y0, x1, x0, newEdgeEquationInt(x1, y1, x0, y0)
	}

	if x0 < x1 {
		return x0, x1, y0, y1, newEdgeEquationInt(x0, y0, x1, y1)
	}
	return x1, x0, y1, y0, newEdgeEquationInt(x1, y1, x0, y0)
}

func drawLineHorizontal(d float64, edge edgeEquation, img draw.Image, x0, x1, y0, y1, incr int) {
	y := y0
	for x := x0; x <= x1; x++ {
		img.Set(x, y, color.White)
		d += edge.a
		if d <= 0.0 {
			d += edge.b
			y += incr
		}
	}
}
