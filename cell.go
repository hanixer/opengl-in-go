package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	rows    = 10
	columns = 10
)

var (
	triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		0.5, -0.5, 0,
		0.5, 0.5, 0,
		-0.5, 0.5, 0,
	}
)

type cell struct {
	drawable  uint32
	alive     bool
	aliveNext bool
	x         int
	y         int
}

type cellsSet [][]*cell

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (c *cell) updateState(cells cellsSet) {

}

func drawAllsells(cells cellsSet, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if (i%2 != 0) == (j%2 == 0) {
				cells[i][j].draw()
			}
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeCells() cellsSet {
	cells := make(cellsSet, rows)
	for r := 0; r < rows; r++ {
		cells[r] = make([]*cell, columns)
		for c := 0; c < columns; c++ {
			cells[r][c] = newCell(c, r)
		}
	}

	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	var xScale float32 = 1.0 / columns
	var yScale float32 = 1.0 / rows
	for i := 0; i < len(square); i += 3 {
		points[i] = (float32(square[i])+0.5+float32(x))*xScale*2.0 - 1.0
		points[i+1] = (float32(square[i+1])+0.5+float32(y))*yScale*2.0 - 1.0
	}
	return &cell{
		drawable: makeVao(points),
		x:        x,
		y:        y,
	}
}
