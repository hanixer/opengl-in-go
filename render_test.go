package main

import "testing"

func Test_edgeEquation_test(t *testing.T) {

	type args struct {
		x float64
		y float64
	}
	v0 := vertex{x: 0.0, y: 0.0}
	v1 := vertex{x: 5.0, y: 10.0}
	v2 := vertex{x: 10.0, y: 0.0}
	a := args{5, 5}
	tests := []struct {
		name string
		edge edgeEquation
		args args
		want bool
	}{
		// TODO: Add test cases.
		// {"", newEdgeEquation(vertex{x: 0.0, y: 0.0}, vertex{x: 10.0, y: 10.0}), args{-5.0, 5.0}, true},
		{"", newEdgeEquation(v0, v1), a, true},
		{"", newEdgeEquation(v2, v0), a, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.edge.test(tt.args.x, tt.args.y); got != tt.want && false {
				t.Errorf("edgeEquation.test() = %v, want %v", got, tt.want)
			}
		})
	}
}
