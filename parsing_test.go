package main

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"reflect"
	"testing"
)

func Test_parseColor(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want color.Color
	}{
		{"", args{"none"}, color.RGBA{}},
		{"", args{"#120000"}, color.RGBA{0x12, 0, 0, 0xFF}},
		{"", args{"#120056"}, color.RGBA{0x12, 0, 0x56, 0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseColor(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSvgString(t *testing.T) {
	d := `
		<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
			<line x1="0" y1="80" x2="100" y2="20" stroke="black" />
		</svg>`

	t.Run("", func(t *testing.T) {
		svg := parseSvgString(d)
		if len(svg.elements) > 0 {
			l := svg.elements[0].(*svgLine)
			assert.Equalf(t, 0.0, l.from.X(), "x1")
			assert.Equalf(t, 80.0, l.from.Y(), "y1")
			assert.Equalf(t, 100.0, l.to.X(), "x2")
			assert.Equalf(t, 20.0, l.to.Y(), "y2")
		} else {
			t.Errorf("no elements")
		}
	})

	d2 := `
		<svg viewBox="0 0 220 100" xmlns="http://www.w3.org/2000/svg" width="123" height="4321">
			<rect width="100" height="100" x="1" y="2" />
		</svg>`
	t.Run("", func(t *testing.T) {
		svg := parseSvgString(d2)
		if len(svg.elements) > 0 {
			r := svg.elements[0].(*svgRect)
			assert.NotNil(t, r)
			assert.Equal(t, 100.0, r.dimension.X())
			assert.Equal(t, 100.0, r.dimension.Y())
			assert.Equal(t, 2.0, r.position.Y())
			assert.Equal(t, 1.0, r.position.X())
			assert.Equal(t, 123.0, svg.width)
			// assert.Equal(t, 4321.0, svg.height)
		} else {
			t.Errorf("no elements")
		}
	})

	d3 := `<svg viewBox="0 0 200 100" xmlns="http://www.w3.org/2000/svg">
	<polygon points="0,100 50,25 50,75.6 100,0" />
  </svg>`
	t.Run("", func(t *testing.T) {
		svg := parseSvgString(d3)
		if len(svg.elements) > 0 {
			pg := svg.elements[0].(*svgPolygon)
			assert.NotNil(t, pg)
			assert.Equal(t, 4, len(pg.points))
			assert.Equal(t, 0.0, pg.points[0].X())
			assert.Equal(t, 100.0, pg.points[0].Y())
			assert.Equal(t, 50.0, pg.points[2].X())
			assert.Equal(t, 75.6, pg.points[2].Y())
			// assert.Equal(t, 4321.0, svg.height)
		} else {
			t.Errorf("no elements")
		}
	})
}
