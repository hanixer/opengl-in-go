package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/png"
	"os"
	"strings"
)

func main() {
	data := `<svg  xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink">
  <rect x="10" y="10" height="100" width="100"
		style="stroke:#ff0000; fill: #0000ff"/>
</svg>`
	r := strings.NewReader(data)
	n := makeXmlDoc(r)
	parseSvg(n)
}

func main22() {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	svvgg := svgGroup{
		elements: []svgElement{
			&svgLine{
				from: mgl64.Vec2{10, 10},
				to:   mgl64.Vec2{70, 80},
			},
			&svgLine{
				from: mgl64.Vec2{10, 33},
				to:   mgl64.Vec2{187, 10},
			},
		},
	}

	drawSvg(img, &svvgg)

	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()

}
