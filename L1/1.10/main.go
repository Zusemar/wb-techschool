package main

import (
	"fmt"
)

func main() {
	a := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	groups := make(map[int][]float32)

	for _, v := range a {
		bin := int(float64(v)/10.0) * 10
		groups[bin] = append(groups[bin], v)
	}

	for k, slice := range groups {
		fmt.Printf("%d: %v\n", k, slice)
	}
}
