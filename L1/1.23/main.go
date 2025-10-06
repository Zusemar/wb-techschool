package main

import (
	"fmt"
)

func pop(slice *[]int, index int) {
	nig := *slice
	sl := append(nig[:index], nig[index+1:]...)
	*slice = sl
}

func main() {
	ms := []int{10, 20, 30, 40}
	fmt.Println(ms)
	pop(&ms, 2)
	fmt.Println(ms)
}
