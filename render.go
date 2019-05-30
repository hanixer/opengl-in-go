package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

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
	frac    float64
	tie     bool
}

func newEdgeEquation(v0, v1 mgl64.Vec2) edgeEquation {
	var eq edgeEquation
	eq.a = v0.Y() - v1.Y()
	eq.b = v1.X() - v0.X()
	eq.c = v0.X()*v1.Y() - v1.X()*v0.Y()
	if eq.a != 0 {
		eq.tie = eq.a > 0
	} else {
		eq.tie = eq.b > 0
	}
	eq.frac = 1
	return eq
}

func newEdgeEquationInt(x0, y0, x1, y1 int) edgeEquation {
	var e edgeEquation
	e.a = float64(y0 - y1)
	e.b = float64(x1 - x0)
	e.c = float64(x0*y1 - x1*y0)
	return e
}

func (edge edgeEquation) test(v mgl64.Vec2) bool {
	w := edge.evaluateP(v)
	return w > 0 || w == 0 && edge.tie
}

func (edge edgeEquation) evaluate(x, y float64) float64 {
	return edge.a*x + edge.b*y + edge.c
}

func (edge edgeEquation) evaluateP(p mgl64.Vec2) float64 {
	return edge.frac * edge.evaluate(p.X(), p.Y())
}

func (edge edgeEquation) evaluateP1(p mgl64.Vec2) float64 {
	return edge.evaluate(p.X(), p.Y())
}

func (edge *edgeEquation) cacheFraction(v mgl64.Vec2) {
	edge.frac = 1.0 / edge.evaluateP(v)
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

func boundingBox(v0, v1, v2 mgl64.Vec2) (x0, y0, x1, y1 float64) {
	x0 = min(min(v0.X(), v1.X()), v2.X())
	y0 = min(min(v0.Y(), v1.Y()), v2.Y())
	x1 = max(max(v0.X(), v1.X()), v2.X())
	y1 = max(max(v0.Y(), v1.Y()), v2.Y())

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

func makeSamples(samples []mgl64.Vec2, x, y float64, n int) {
	frac := 1.0 / float64(n)
	if n == 1 {
		samples[0][0] = x + 0.5
		samples[0][1] = y + 0.5
		return
	}
	for i := 0.0; i < float64(n); i++ {
		for j := 0.0; j < float64(n); j++ {
			xx := x + j*frac + rand.Float64()*frac
			yy := y + i*frac + rand.Float64()*frac
			samples[int(i)*n+int(j)][0] = xx
			samples[int(i)*n+int(j)][1] = yy
		}
	}
}

func fillTriangle(v0, v1, v2 mgl64.Vec2, img *image.RGBA, fillColor color.RGBA, samplesCount int) {
	minX, minY, maxX, maxY := boundingBox(v0, v1, v2)

	e0 := newEdgeEquation(v1, v2)
	e1 := newEdgeEquation(v2, v0)
	e2 := newEdgeEquation(v0, v1)

	e0.cacheFraction(v0)
	e1.cacheFraction(v1)
	e2.cacheFraction(v2)

	samples := make([]mgl64.Vec2, samplesCount*samplesCount)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			var colorAcc [4]uint64
			makeSamples(samples, x, y, samplesCount)
			for _, p := range samples {
				// compute baricentric coordinates
				if e0.test(p) && e1.test(p) && e2.test(p) {
					addColor(&colorAcc, fillColor)
				}
			}
			divideColor(&colorAcc, samplesCount*samplesCount)
			if colorAcc[3] != 0 {
				c64 := color.RGBA{convertColor64(colorAcc[0]), convertColor64(colorAcc[1]),
					convertColor64(colorAcc[2]), convertColor64(colorAcc[3])}
				// img.Set(int(x), int(y), c64)
				setPixel(img, x, y, c64)
			}
		}

	}
}

func convertColor64(v uint64) uint8 {
	return uint8(v * 0xFF / 0xFFFF)
}

func setPixel(img *image.RGBA, x, y float64, fillColor color.RGBA) {
	// Ca' = 1 - (1 - Ea) * (1 - Ca)
	// Cr' = (1 - Ea) * Cr + Er
	// Cg' = (1 - Ea) * Cg + Eg
	// Cb' = (1 - Ea) * Cb + Eb
	xx := int(x)
	yy := int(y)
	// max := 0xFFFF
	// fr, fg, fb, _ := fillColor.R
	// br, bg, bb, ba := img.At(xx, yy).RGBA()
	// r := fillColor.R + img
	img.Set(xx, yy, fillColor)
	// img.Set(xx, yy, color.RGBA64{uint16((fr + br) * 0xFFFF / 0xFFFFFFFF), uint16(fg + bg), uint16(fb + bb), uint16(ba)})
	// cr = (max-ea)*cr + er
	// cg = (max-ea)*cg + eg
	// cb = (max-ea)*cb + eb
	// ca = max - (max-ea)*(max-ca)
	// img.Set(xx, yy, color.RGBA64{})

}

func addColor(colorAccum *[4]uint64, c color.Color) {
	r, g, b, a := c.RGBA()
	colorAccum[0] += uint64(r)
	colorAccum[1] += uint64(g)
	colorAccum[2] += uint64(b)
	colorAccum[3] += uint64(a)
}

func divideColor(colorAccum *[4]uint64, samplesCount int) {
	colorAccum[0] /= uint64(samplesCount)
	colorAccum[1] /= uint64(samplesCount)
	colorAccum[2] /= uint64(samplesCount)
	colorAccum[3] /= uint64(samplesCount)
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
