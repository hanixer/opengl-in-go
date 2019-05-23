package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var defaultInputSvg = "svg/line.svg"

func main() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	n := parseSvgString(string(data))
	img := drawSvg(n)
	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}
