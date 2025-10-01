package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var justString string

func main() {
	fmt.Println("=== До исправления ===")
	someFunc()

	fmt.Println("\n=== После исправления ===")
	goodFunc()
}

func createHugeString(n int) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString("Hello")
	}
	return builder.String()
}

// этот способ создание строки justString приведет к утечке памяти
// потому что justString будет хранить указатель на весь v а не
// только на первые 100 символов и когда строка v уже будет не нужна
// сборщик все равно ее не очистит потому что justString на нее ссылается
func someFunc() {
	v := createHugeString(1 << 10) // примерно 5 МБ (1024*1024*5 байт)
	justString = v[:100]

	// Проверяем адреса структур
	vh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	jh := (*reflect.StringHeader)(unsafe.Pointer(&justString))
	fmt.Printf("v struct addr: %p, v.data addr: %x\n", &v, vh.Data)
	fmt.Printf("justString struct addr: %p, justString.data addr: %x\n", &justString, jh.Data)
}

// тут напрямую создаем через copy обособленную структуру b, у которой свой кусок памяти,
// не ссылающийся на слайс даты, закрепленный за v
func goodFunc() {
	v := createHugeString(1 << 20) // тоже ~5 МБ
	// копируем только первые 100 символов
	b := make([]byte, 100)
	copy(b, v[:100])
	justString = string(b)

	vh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	jh := (*reflect.StringHeader)(unsafe.Pointer(&justString))
	fmt.Printf("v struct addr: %p, v.data addr: %x\n", &v, vh.Data)
	fmt.Printf("justString struct addr: %p, justString.data addr: %x\n", &justString, jh.Data)
}
