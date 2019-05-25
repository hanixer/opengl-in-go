package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

var defaultInputSvg = `svg\polygons2.svg`

func main() {
	// data, err := ioutil.ReadFile(defaultInputSvg)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// n := parseSvgString(string(data))
	// img := drawSvg(n)
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)
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
	// drawTriangle(vert(0, 0), vert(2, 2), vert(0, 2), img, color.White)
	drawLine(0, 50, 500, 40, img, color.RGBA{255, 0, 0, 255})

	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}
