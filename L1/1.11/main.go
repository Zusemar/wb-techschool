package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	b := []int{5, 6, 7, 8, 9, 10, 11, 12, 13}

	c := make(map[int]int)

	for _, val := range a {
		c[val]++
	}

	for _, val := range b {
		c[val]++
	}

	fmt.Println(c)
	ans := []int{}
	for key, val := range c {
		if val == 2 {
			ans = append(ans, key)
		}
	}
	sort.Ints(ans)
	fmt.Println(ans)
}
