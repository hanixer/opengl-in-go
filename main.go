package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var defaultInputSvg = `svg\illustration\06_sphere.svg`
var samplesCount = 1
var texture uint32

func main5() {
	// data, err := ioutil.ReadFile(defaultInputSvg)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// n := parseSvgString(string(data))
	// img := drawSvg(n)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

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
	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	vao := makeVao(vertexData)
	texture = makeTexture()

	updateSvg()
	// window.SetSize(int(svgData.width), int(svgData.height))

	gl.UseProgram(program)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		// draw(window, prog)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func updateSvg() {
	data, err := ioutil.ReadFile(defaultInputSvg)
	if err != nil {
		log.Fatalln(err)
	}
	svgData := parseSvgString(string(data))
	img := drawSvg(svgData, samplesCount)
	swapVertically(img)
	updateTexture(texture, img)
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
		}
	}
}
