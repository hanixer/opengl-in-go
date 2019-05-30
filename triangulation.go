package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

// listNode is a part of double linked list of points
type listNode struct {
	prev, next *listNode
	point      mgl64.Vec2
}

// make a double linked list of point from a slice
func constructList(points []mgl64.Vec2, reverse bool) *listNode {
	var prev, head *listNode
	for _, point := range points {
		curr := new(listNode)
		curr.point = point
		curr.prev = prev
		if prev != nil {
			prev.next = curr
		}
		prev = curr
		if head == nil {
			head = curr
		}
	}
	prev.next = head
	head.prev = prev
	if reverse {
		node := head
		for i := 0; i < len(points); i++ {
			next := node.next
			node.next = node.prev
			node.prev = next
			node = next
		}
	}
	return head
}

// triangulate takes points of a polygon and returns triangles,
// that cover the polygon.
// implementation uses ear cutting method
func triangulate(points []mgl64.Vec2) []mgl64.Vec2 {
	result := []mgl64.Vec2{}
	reverse := shouldReverse(points)
	node := constructList(points, reverse)
	listSize := len(points)

	for listSize > 2 {
		tip := findEar(node, listSize)
		if tip != nil {
			result = append(result, tip.prev.point)
			result = append(result, tip.point)
			result = append(result, tip.next.point)

			tip.prev.next = tip.next
			tip.next.prev = tip.prev
			listSize--
			node = tip.next
		} else {
			return result
		}
	}

	return result
}

// shouldReverse returns true if we need to reverse order of points
// before triangulation
func shouldReverse(points []mgl64.Vec2) bool {
	point := points[0]
	index := 0
	for i := 1; i < len(points); i++ {
		if points[i].X() < point.X() && points[i].Y() < point.Y() {
			point = points[i]
			index = i
		}
	}
	var p0, p1, p2 mgl64.Vec2
	if index == 0 {
		p0 = points[len(points)-1]
	} else {
		p0 = points[index-1]
	}
	p1 = points[index]
	p2 = points[(index+1)%len(points)]
	return !isConvexAngle(p0, p1, p2)
}

// findEar returns a "tip of ear", i.e. if points p0, p1, p2 forms a "ear" triangle
// return the node containing point p1
func findEar(node *listNode, listSize int) *listNode {
	for i := 0; i < listSize; i++ {
		p0 := node.prev.point
		p1 := node.point
		p2 := node.next.point
		if isConvexAngle(p0, p1, p2) {
			if !anyPointInsideTriangle(p0, p1, p2, node, listSize) {
				return node
			}
		}
		node = node.next
	}
	return nil
}

// isConvexAngle returns true if given points form an angle which is less that pi radians
func isConvexAngle(p0, p1, p2 mgl64.Vec2) bool {
	v0 := p0.Sub(p1)
	v1 := p2.Sub(p1)
	z := v0.X()*v1.Y() - v1.X()*v0.Y()
	return z > 0
}

func anyPointInsideTriangle(p0, p1, p2 mgl64.Vec2, head *listNode, listSize int) bool {
	node := head
	for i := 0; i < listSize; i++ {
		if node.point != p0 && node.point != p1 && node.point != p2 {
			if isPointInsideTriangle(p0, p1, p2, node.point) {
				return true
			}
		}
		node = node.next
	}
	return false
}

func isPointInsideTriangle(p0, p1, p2, point mgl64.Vec2) bool {
	e0 := newEdgeEquation(p1, p2)
	e1 := newEdgeEquation(p2, p0)
	e2 := newEdgeEquation(p0, p1)

	w0 := e0.evaluate(point.X(), point.Y())
	w1 := e1.evaluate(point.X(), point.Y())
	w2 := e2.evaluate(point.X(), point.Y())

	return w0 <= 0 && w1 <= 0 && w2 <= 0
}
