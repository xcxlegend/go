package main

import (
	"fmt"
)

/**
 * 快速排序
 * @param {[type]} array []int [description]
 */
func QuickSort(array []int) {
	quick_sort(array, 0, len(array)-1)
}

func quick_sort(array []int, start, end int) {
	if start >= end {
		return
	}

	var l = start
	var r = end

	var temp = array[start]
	for l < r {

		for l < r {

			// 小数往左排
			if array[r] < temp {
				array[l] = array[r]
				l++
				break
			}
			r--
		}

		for l < r {
			// ! 关键 大数往右排
			if array[l] > temp {
				array[r] = array[l]
				r--
				break
			}
			l++
		}

	}
	array[l] = temp
	quick_sort(array, start, l-1)
	quick_sort(array, l+1, end)

}

func main() {
	var array = []int{1, 3, 4, 75, 65, 32, 12, 35, 77, 23}
	QuickSort(array)
	fmt.Println(array)
}
