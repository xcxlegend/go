package main

import (
	"fmt"
	"math"
)

func findNearestSortScores(arr []int, k, x int) []int{
	i := 0
	for i < len(arr) - k{
		if math.Abs(float64(arr[i] - x)) > math.Abs(float64(arr[i + k] - x)){
			i++
		}else{
			break
		}
	}
	return arr[i:i+k]
}


func main() {
	arr := []int{0, 0, 1, 2, 3, 3, 4, 6, 7, 8}
	k := 4
	x := 3
	n := findNearestSortScores(arr, k, x)
	fmt.Println(n)
}
