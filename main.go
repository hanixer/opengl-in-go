package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var defaultInputSvg = `svg\subdiv\triangle1.svg`

func main5() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	n := parseSvgString(string(data))
	img := drawSvg(n)
	// img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// numSamples := 10
	// samples := make([]mgl64.Vec2, numSamples*numSamples)
	// makeSamples(samples, 0, 0, numSamples)
	// for _, v := range samples {
	// 	i, j := int(v.X()*100), int(v.Y()*100)
	// 	fmt.Println(v)
	// 	img.Set(i, j, color.Black)
	// }

	// drawLine(0, 50, 500, 40, img, color.RGBA{255, 0, 0, 255})
	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}
