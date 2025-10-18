package main

import "fmt"

func main() {
	fmt.Println(stringInverter("sinep"))
}

func stringInverter(str string) string {
	runes := []rune(str)
	reversed := make([]rune, len(runes))
	for i, r := range runes {
		reversed[len(runes)-1-i] = r
	}
	return string(reversed)
}
