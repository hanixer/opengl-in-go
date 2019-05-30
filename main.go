package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var defaultInputSvg = `svg\hexes.svg`
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

	img2 := image.NewRGBA(image.Rect(0, 0, int(svgData.width), int(svgData.height)))
	valu := 100.0
	fillTriangle(mgl64.Vec2{0, 0}, mgl64.Vec2{0, valu}, mgl64.Vec2{valu, 0}, img2, color.White, 2)
	fillTriangle(mgl64.Vec2{0, valu}, mgl64.Vec2{valu, valu}, mgl64.Vec2{valu, 0}, img2, color.RGBA{0xFF, 0, 0, 0xFF}, 2)

	f, _ := os.Create("out.png")
	png.Encode(f, img2)
	f, _ = os.Create("out1.png")
	png.Encode(f, img)

	swapVertically(img)

	updateTexture(texture, img)
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
