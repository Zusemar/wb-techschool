package main

import "fmt"

// https://habr.com/ru/companies/ruvds/articles/744230/

func main() {
	var n int64
	fmt.Println("\nВведи число")
	fmt.Scanf("%d\n", &n)
	fmt.Printf("%064b", n)

	var i int64
	fmt.Println("\nВведи бит")
	fmt.Scanf("%d\n", &i)

	var val int64
	fmt.Println("\nВведи значение i бита")
	fmt.Scanf("%d\n", &val)

	fmt.Printf("%064b", bitSetter(n, i, val))
}

func bitSetter(n int64, i int64, val int64) int64 {
	if val == 0 {
		n = n & ^(1 << i)
	} else {
		return n | (1 << i)
	}
	return n
}
