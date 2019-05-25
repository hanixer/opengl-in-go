package main

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stretchr/testify/assert"
)

func Test_constructList(t *testing.T) {
	t.Run("", func(t *testing.T) {
		p1, p2, p3 := mgl64.Vec2{1, 1}, mgl64.Vec2{2, 2}, mgl64.Vec2{3, 3}
		points := []mgl64.Vec2{p1, p2, p3}
		res := constructList(points)
		assert.Equal(t, p1, res.point)
		assert.Equal(t, p2, res.next.point)
		assert.Equal(t, p3, res.next.next.point)
		assert.Equal(t, p1, res.next.next.next.point)
		assert.Equal(t, p3, res.prev.point)
		assert.Equal(t, p2, res.prev.prev.point)
		assert.Equal(t, p1, res.prev.prev.prev.point)
	})
}

func Test_isPointInsideTriangle(t *testing.T) {
	type args struct {
		p0    mgl64.Vec2
		p1    mgl64.Vec2
		p2    mgl64.Vec2
		point mgl64.Vec2
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{mgl64.Vec2{0, 0}, mgl64.Vec2{0, 5}, mgl64.Vec2{5, 5}, mgl64.Vec2{3, 3}}, true},
		{"", args{mgl64.Vec2{0, 0}, mgl64.Vec2{0, 5}, mgl64.Vec2{5, 5}, mgl64.Vec2{-1, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPointInsideTriangle(tt.args.p0, tt.args.p1, tt.args.p2, tt.args.point); got != tt.want {
				t.Errorf("isPointInsideTriangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_triangulate(t *testing.T) {
	p0 := mgl64.Vec2{3, 3}
	p1 := mgl64.Vec2{6, 7}
	p2 := mgl64.Vec2{7, 9}
	p3 := mgl64.Vec2{9, 3}
	p4 := mgl64.Vec2{6, 5}
	p5 := mgl64.Vec2{5, 5}
	p6 := mgl64.Vec2{5, 10}
	p7 := mgl64.Vec2{10, 10}
	p8 := mgl64.Vec2{10, 5}
	q0 := mgl64.Vec2{0, 0}
	q1 := mgl64.Vec2{0, 5}
	q2 := mgl64.Vec2{5, 5}
	q3 := mgl64.Vec2{5, 4}
	q4 := mgl64.Vec2{1, 4}
	q5 := mgl64.Vec2{1, 0}
	type args struct {
		points []mgl64.Vec2
	}
	tests := []struct {
		name string
		args args
		want []mgl64.Vec2
	}{
		// {"", args{[]mgl64.Vec2{p0, p1, p4}}, []mgl64.Vec2{p0, p1, p4}},
		{"", args{[]mgl64.Vec2{q0, q1, q2, q3, q4, q5}}, []mgl64.Vec2{q5, q0, q1, q1, q2, q3, q1, q3, q4, q1, q4, q5}},
		{"", args{[]mgl64.Vec2{p5, p6, p7, p8}}, []mgl64.Vec2{p8, p5, p6, p8, p6, p7}},
		{"", args{[]mgl64.Vec2{p0, p1, p2, p3, p4}}, []mgl64.Vec2{p4, p0, p1, p4, p1, p2, p4, p2, p3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := triangulate(tt.args.points); !isSameTriangles(got, tt.want) {
				t.Errorf("triangulate() = %v, want %v", got, tt.want)
			}
		})
	}
}

type triangles [][]mgl64.Vec2

func (t triangles) contains(t0 []mgl64.Vec2) bool {
	for _, t1 := range t {
		if isSameTri(t0, t1) {
			return true
		}
	}
	return false
}

func isSameTriangles(t0, t1 []mgl64.Vec2) bool {
	if len(t0) != len(t1) {
		return false
	}

	triangles0 := make(triangles, 0)
	triangles1 := make(triangles, 0)

	for i := 0; i < len(t0); i += 3 {
		triangles0 = append(triangles0, t0[i:i+3])
		triangles1 = append(triangles1, t1[i:i+3])
	}

	for _, tt0 := range triangles0 {
		if !triangles1.contains(tt0) {
			return false
		}
	}

	return true
}

func isSameTri(t0, t1 []mgl64.Vec2) bool {
	if len(t0) != len(t1) {
		return false
	}

	p := t0[0]
	for i := 0; i < len(t1); i++ {
		if t1[i] == p {
			for j := 1; j < len(t0); j++ {
				i = (i + 1) % len(t1)
				if t0[j] != t1[i] {
					return false
				}
			}
			return true
		}
	}
	return false
}
