package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var defaultInputSvg = `svg\hardcore\01_degenerate_square1.svg`

func main() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	n := parseSvgString(string(data))
	img := drawSvg(n)
	// img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	// drawLine(0, 50, 500, 40, img, color.RGBA{255, 0, 0, 255})
	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}
