package main

import "fmt"

func main() {
	typeReader(42)
	typeReader("hello")
	typeReader(true)
	typeReader(make(chan int))
}

func typeReader(mf interface{}) {
	switch v := mf.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	case bool:
		fmt.Println("bool:", v)
	case chan int:
		fmt.Println("chan int")
	default:
		fmt.Println("unknown type")
	}
}
