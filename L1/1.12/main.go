package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []string{"cat", "cat", "dog", "cat", "tree", "tree"}

	b := make(map[string]int)

	ans := []string{}

	for _, val := range a {
		b[val]++
	}

	for key, _ := range b {
		ans = append(ans, key)
	}

	sort.Strings(ans)

	fmt.Println(ans)
}
