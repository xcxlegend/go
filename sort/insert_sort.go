package main

import (
	"fmt"
)

/**
 * 插入排序
 * @param {[type]} array []int [description]
 */
func InsertSort(array []int) {

	// for i := 1; i < len(array); i++ {
	for i, value := range array {
		var flag = i
		var j = i - 1
		for ; j > 0; j-- {

			if value < array[j] {
				flag = j
			}
		}

		// for j = i - 1; j >= flag; j-- {
		// 	array[j+1] = array[j]
		// }

		copy(array[flag+1:i+1], array[flag:i])

		array[flag] = value

	}

}

func main() {
	var array = []int{1, 3, 99, 4, 75, 65, 32, 12, 35, 77, 23}

	InsertSort(array)
	fmt.Println(array)
}
