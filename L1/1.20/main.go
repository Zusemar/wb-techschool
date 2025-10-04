package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "snow dog sun"
	var builder strings.Builder

	// распиливаем строчку на массив слов (проигнорим strings.Field)
	var words []string
	var builder1 strings.Builder
	for _, r := range str {
		if r == ' ' {
			if builder1.Len() > 0 {
				words = append(words, builder1.String())
				builder1.Reset()
			}
		} else {
			builder1.WriteRune(r)
		}
	}
	// добавим последнее слово
	if builder1.Len() > 0 {
		words = append(words, builder1.String())
	}

	// собираем собранные слова
	var ans string
	for i := len(words) - 1; i >= 0; i-- {
		builder.WriteString(words[i])
		if i > 0 {
			builder.WriteString(" ")
		}
	}
	ans = builder.String()
	fmt.Println(ans)
}
