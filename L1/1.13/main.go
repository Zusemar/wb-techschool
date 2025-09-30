package main

import "fmt"

func main() {
	a, b := 5, 10

	// Обмен через сложение/вычитание
	a = a + b
	b = a - b
	a = a - b
	fmt.Println("After swap (sum/diff):", a, b)

	// Обмен через XOR
	a, b = 5, 10
	a = a ^ b
	b = a ^ b
	a = a ^ b
	fmt.Println("After swap (XOR):", a, b)
}
