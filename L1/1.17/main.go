package main

import "fmt"

func main() {
	fmt.Println(binSearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 9))
}

func binSearch(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		mid := l + (r-l)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return -1
}
