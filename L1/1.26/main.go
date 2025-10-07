package main

import (
	"fmt"
	"strings"
)

func main() {
	str1 := "aabbccdd"
	str2 := "abcdABCD"
	str3 := "abcdEFG"

	fmt.Println(isRepetitions(str1), "\n", isRepetitions(str2), "\n", isRepetitions(str3))
}

func isRepetitions(str string) bool {
	strLower := strings.ToLower(str)
	seen := make(map[rune]bool)
	for _, ch := range strLower {
		if seen[ch] {
			return false
		}
		seen[ch] = true
	}
	return true
}
