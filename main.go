package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/png"
	"os"
)

func main() {
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

	xml.Unmarshal

	drawSvg(img, &svvgg)

	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()

}
