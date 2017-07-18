package main

import (
	"fmt"
)

/**
 * 冒泡排序
 * @param {[type]} array []int [description]
 */
func BubbleSort(array []int) {

	var count = len(array)
	for i := 0; i < count; i++ {

		for j := i; j < count-1; j++ {

			if array[j] > array[j+1] {
				array[j+1], array[j] = array[j], array[j+1]
			}

		}

	}

}

func main() {
	var array = []int{1, 3, 4, 75, 65, 32, 12, 35, 77, 23}

	BubbleSort(array)
	fmt.Println(array)
}
