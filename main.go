package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width              = 500
	height             = 500
	vertexShaderSource = `
#version 410
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;

out vec2 texCoord;

void main() {
	gl_Position = vec4(aPos, 1.0);
	texCoord = aTexCoord;
}
` + "\x00"

	fragmentShaderSource = `
#version 410
in vec2 texCoord;

out vec4 fragColor;

uniform sampler2D ourTexture;
void main() {
	fragColor = texture(ourTexture, texCoord);
	// fragColor = vec4(0.3, 0.3, 0.5, 1);
}
` + "\x00"
)

var (
	vertexData = []float32{
		1.0, 1.0, 0.0, 1.0, 1.0, // top right
		1.0, -1.0, 0.0, 1.0, 0.0, // bottom right
		-1.0, -1.0, 0.0, 0.0, 0.0, // bottom left
		-1.0, 1.0, 0.0, 0.0, 1.0, // top left
	}
	indices = []uint32{
		0, 1, 3,
		1, 2, 3,
	}
)

func main2() {
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	vao := makeVao(vertexData)
	texture := makeTexture()

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

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	v0 := vertex{x: 10, y: 10}
	v1 := vertex{x: 200, y: 100}
	v2 := vertex{x: 350, y: 10}
	drawTriangle(v0, v1, v2, img, color.Black)
	f, _ := os.Create("out.png")
	png.Encode(f, img)
	f.Close()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Opengl thing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalln(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func draw123(window *glfw.Window, program uint32, vao uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(0))

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeTexture() uint32 {
	f, _ := os.Open("1.png")
	img2, _ := png.Decode(f)

	img := image.NewRGBA(img2.Bounds())
	draw.Draw(img, img.Bounds(), img2, image.ZP, draw.Src)

	height := img.Bounds().Dy()
	stride := make([]uint8, img.Stride)
	for row := 0; row < height/2; row++ {
		mirrorRow := height - row - 1
		strideCurrent := img.Pix[row*img.Stride : row*img.Stride+img.Stride]
		strideMirror := img.Pix[mirrorRow*img.Stride : mirrorRow*img.Stride+img.Stride]
		copy(stride, strideCurrent)
		copy(strideCurrent, strideMirror)
		copy(strideMirror, stride)
	}

	c := image.NewUniform(color.RGBA{30, 100, 40, 255})
	for i := 0; i < 20; i += 2 {
		r := image.Rect(30, i+10, 60, i+20)
		draw.Draw(img, r, c, image.ZP, draw.Src)
	}

	ptr := &img.Pix[0]
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(ptr))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.CreateBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var ebo uint32
	gl.CreateBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csourses, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csourses, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
