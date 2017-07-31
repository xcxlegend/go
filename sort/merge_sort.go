package main

import (
	"fmt"
)

/**
 * 归并排序
 * @param {[type]} array []int [description]
 */
func MergeSort(array []int) {

	var count = len(array)
	if count <= 1 {
		return
	}

	var temp = make([]int, count)
	copy(temp, array)
	var half = (count - 1) / 2
	merge_sort(array, temp, 0, half)
	merge_sort(array, temp, half+1, count-1)
	merge(array, temp, 0, half, count-1)
}

func merge_sort(array, temp []int, l, r int) {

	if l >= r {
		return
	}

	var half = (l + r) / 2

	// ! 归并排序重点在于 `分` `合`
	merge_sort(array, temp, l, half)
	merge_sort(array, temp, half+1, r)
	merge(array, temp, l, half, r)
}

func merge(array, temp []int, l, half, r int) {

	var p1 = l
	var p2 = half + 1
	var i = l
	for p1 <= half && p2 <= r {
		if array[p1] <= array[p2] {
			temp[i] = array[p1]
			p1++
		} else {
			temp[i] = array[p2]
			p2++
		}
		i++
	}

	if p2 > r {
		copy(temp[i:r+1], array[p1:half+1])
	} else {
		copy(temp[i:r+1], array[p2:r+1])
	}

	copy(array[l:r+1], temp[l:r+1])
}

func main() {
	var array = []int{1, 3, 99, 4, 75, 65, 32, 12, 35, 77, 23}
	MergeSort(array)
	fmt.Println(array)
}
