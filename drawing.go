package main

import (
	"image"
	"image/draw"
)

func drawSvg(svg *svg) draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(svg.width), int(svg.height)))

	drawElements(svg.elements, img)

	return img
}

func drawElements(elements []svgElement, img draw.Image) {
	for _, elem := range elements {
		switch v := elem.(type) {
		case *svgLine:
			drawLineF(v.from.X(), v.from.Y(), v.to.X(), v.to.Y(), img)
		}
	}
}
