package main

import (
	"fmt"
)

func main() {
	fmt.Println(qsort([]int{5, 9, 3, 0, 1, 7, 12}))
}

// Выбираем опорный элемент (pivot).
// Разделяем массив на две части:
// элементы меньше pivot и больше (или равные) pivot.
// Рекурсивно сортируем обе части.
func qsort(a []int) []int {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivotIndex := left

	a[pivotIndex], a[right] = a[right], a[pivotIndex]

	for i := range a {
		if a[i] < a[right] {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]

	qsort(a[:left])
	qsort(a[left+1:])

	return a
}
