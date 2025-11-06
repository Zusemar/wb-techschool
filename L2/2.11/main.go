package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagrams(words))
}

func anagrams(words []string) map[string][]string {
	groups := make(map[string][]string)
	order := make(map[string]string) // запоминает первое слово для ключа

	for _, w := range words {
		w = strings.ToLower(w)
		runes := []rune(w)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)

		if _, ok := order[key]; !ok {
			order[key] = w
		}
		groups[key] = append(groups[key], w)
	}

	result := make(map[string][]string)
	for key, group := range groups {
		if len(group) > 1 {
			sort.Strings(group)
			result[order[key]] = group
		}
	}

	return result
}
