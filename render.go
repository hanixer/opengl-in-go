package main

import (
	"image/color"
	"image/draw"
	"math"

	"github.com/go-gl/mathgl/mgl64"
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

func newEdgeEquation2(v0, v1 mgl64.Vec2) edgeEquation {
	var eq edgeEquation
	eq.a = v0.Y() - v1.Y()
	eq.b = v1.X() - v0.X()
	eq.c = v0.X()*v1.Y() - v1.X()*v0.Y()
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

	x0 = math.Floor(x0)
	y0 = math.Floor(y0)
	x1 = math.Ceil(x1)
	y1 = math.Ceil(y1)

	return
}

func makeColor(r, g, b float64) color.Color {
	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

func vert(x, y float64) vertex { return vertex{x: x, y: y} }

func convertVertex(v mgl64.Vec2) vertex { return vertex{x: v.X(), y: v.Y()} }

func fillTriangle2(v0, v1, v2 mgl64.Vec2, img draw.Image, c color.Color) {
	fillTriangle(convertVertex(v0), convertVertex(v1), convertVertex(v2), img, c)
}

func fillTriangle(v0, v1, v2 vertex, img draw.Image, c color.Color) {
	minX, minY, maxX, maxY := boundingBox(v0, v1, v2)

	e0 := newEdgeEquation(v1, v2)
	e1 := newEdgeEquation(v2, v0)
	e2 := newEdgeEquation(v0, v1)

	k0 := 1.0 / e0.evaluate(v0.x, v0.y)
	k1 := 1.0 / e1.evaluate(v1.x, v1.y)
	k2 := 1.0 / e2.evaluate(v2.x, v2.y)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// compute baricentric coordinates
			w0 := k0 * e0.evaluate(x, y)
			w1 := k1 * e1.evaluate(x, y)
			w2 := k2 * e2.evaluate(x, y)

			// fmt.Printf("x = %v; y = %v; k0 = %v; w0 = %v; w1 = %v; w2 = %v\n", x, y, k0, w0, w1, w2)
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				img.Set(int(x), int(y), c)
			}
		}
	}
}

func drawLinePoints(p0, p1 mgl64.Vec2, img draw.Image, fillColor color.Color) {
	drawLineF(p0.X(), p0.Y(), p1.X(), p1.Y(), img, fillColor)
}
func drawLineF(x0, y0, x1, y1 float64, img draw.Image, fillColor color.Color) {
	drawLine(int(x0), int(y0), int(x1), int(y1), img, fillColor)
}

func drawPoint(x0, y0 float64, img draw.Image, fillColor color.Color) {
	img.Set(int(x0), int(y0), fillColor)
}

func drawLine(x0, y0, x1, y1 int, img draw.Image, fillColor color.Color) {
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
			img.Set(v, u, fillColor)
			if incr > 0 && edge.evaluate(float64(v+incr), float64(u)+1.5) > 0 {
				v += incr
			} else if incr < 0 && edge.evaluate(float64(v+incr), float64(u)+1.5) < 0 {
				v += incr
			}
		} else {
			img.Set(u, v, fillColor)
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
