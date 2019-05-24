package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var defaultInputSvg = `svg\polygons.svg`

func main() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	n := parseSvgString(string(data))
	img := drawSvg(n)
	img.At(1, 1)
	// img := image.NewRGBA(image.Rect(0, 0, 51, 51))
	// draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)
	// drawLine(0, 50, 50, 50, img)
	// drawLine(0, 0, 10, 50, img)
	// drawLine(0, 0, 50, 30, img)
	// drawLine(0, 0, 50, 0, img)
	// drawLine(0, 0, 0, 50, img)
	// drawLine(0, 50, 50, 40, img)
	// drawLine(0, 50, 10, 0, img)
	// drawLine(0, 50, 50, 50, img)
	// drawLine(0, 50, 0, 0, img)
	// drawLine(20, 20, 20, 15, img)
	// drawTriangle(vert(0, 0), vert(0, 2), vert(2, 2), img, color.Black)

	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}
