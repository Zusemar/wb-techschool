package main

import (
	"fmt"
	"math"
)

type Point struct {
	x float64
	y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func (p *Point) Distance(other *Point) float64 {
	return math.Hypot(p.x-other.x, p.y-other.y)
}

func main() {

	p1 := NewPoint(0, 0)
	p2 := NewPoint(1, 1)
	fmt.Println(p1.Distance(p2))
}
