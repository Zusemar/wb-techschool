package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Задаём большие числа
	a := big.NewInt(0)
	b := big.NewInt(0)
	a.SetString("1048577", 10) // чуть больше 2^20
	b.SetString("2097153", 10) // примерно 2^21

	// Результаты операций
	sum := new(big.Int).Add(a, b)
	diff := new(big.Int).Sub(a, b)
	prod := new(big.Int).Mul(a, b)
	quot := new(big.Int).Div(a, b)

	fmt.Println("a + b =", sum)
	fmt.Println("a - b =", diff)
	fmt.Println("a * b =", prod)
	fmt.Println("a / b =", quot)
}
