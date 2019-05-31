package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var defaultInputSvg = `svg\illustration\02_hexes.svg`
var samplesCount = 1
var texture uint32
var window *glfw.Window
var svgPaths []string
var svgIndex int
var loadOnlyDefault = true

func walkPath(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(path, ".svg") {
		svgPaths = append(svgPaths, path)
	}
	return nil
}

func main55() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	n := parseSvgString(string(data))
	img := drawSvg(n, samplesCount)
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

func main() {
	runtime.LockOSThread()

	filepath.Walk("svg", walkPath)

	window = initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	vao := makeVao(vertexData)
	texture = makeTexture()

	updateSvg()

	gl.UseProgram(program)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		glfw.PollEvents()
		window.SwapBuffers()
		window.SetSizeCallback(sizeCallback)
	}
}

func updateSvg() {
	var svgPath string
	if loadOnlyDefault {
		svgPath = defaultInputSvg
	} else {
		svgPath = svgPaths[svgIndex]
	}

	window.SetTitle(svgPath)
	data, err := ioutil.ReadFile(svgPath)
	if err != nil {
		log.Fatalln(err)
	}
	svgData := parseSvgString(string(data))
	window.SetSize(int(svgData.width), int(svgData.height))
	gl.Viewport(0, 0, int32(svgData.width), int32(svgData.height))
	img := drawSvg(svgData, samplesCount)
	// img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	var f *os.File
	f, _ = os.Create("out1.png")
	png.Encode(f, img)

	img2 := image.NewRGBA(image.Rect(0, 0, 500, 500))
	color1 := color.RGBA{255, 129, 104, 255}
	color2 := color.RGBA{23, 249, 237, 20}
	color3 := mixColors(color1, color2)
	draw.Draw(img2, image.Rect(0, 0, 100, 100), image.NewUniform(color1), image.ZP, draw.Src)
	draw.Draw(img2, image.Rect(100, 0, 200, 100), image.NewUniform(color2), image.ZP, draw.Src)
	draw.Draw(img2, image.Rect(200, 0, 300, 100), image.NewUniform(color3), image.ZP, draw.Src)

	f, _ = os.Create("out.png")
	png.Encode(f, img2)

	swapVertically(img)

	updateTexture(texture, img)
}

func mixColors(c1, c2 color.RGBA) color.RGBA {
	al := float64(c2.A) / 255.0
	r := uint8(float64(c2.R)*al + float64(c1.R)*(1.0-al))
	g := uint8(float64(c2.G)*al + float64(c1.G)*(1.0-al))
	b := uint8(float64(c2.B)*al + float64(c1.B)*(1.0-al))
	return color.RGBA{r, g, b, 255}
}

func sizeCallback(_ *glfw.Window, w int, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))
}

func keyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release {
		if key == glfw.KeyMinus {
			if samplesCount > 1 {
				samplesCount--
				updateSvg()
			}
		} else if key == glfw.KeyBackspace {
			samplesCount++
			fmt.Println(samplesCount)
			updateSvg()
		} else if key == glfw.KeyR {
			w, h := window.GetSize()
			window.SetSize(w+100, h+100)
		} else if key == glfw.KeyS {
			w, h := window.GetSize()
			window.SetSize(w-100, h-100)
		} else if key == glfw.KeyLeft {
			svgIndex--
			if svgIndex < 0 {
				svgIndex = len(svgPaths) - 1
			}
			updateSvg()
		} else if key == glfw.KeyRight {
			svgIndex = (svgIndex + 1) % len(svgPaths)
			updateSvg()
		}
	}
}
