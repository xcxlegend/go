package main

import (
	"fmt"
)

/**
 * 选择排序
 * @param {[type]} array []int [description]
 */
func ChooseSort(array []int) {
	var count = len(array)
	for i := 0; i < count; i++ {

		for j := i; j < count; j++ {

			// 将冒泡排序的判断稍微改一下即可
			//
			if array[i] > array[j] {
				array[j], array[i] = array[i], array[j]
			}

		}

	}

}

func main() {
	var array = []int{1, 3, 4, 75, 65, 32, 12, 35, 77, 23}

	ChooseSort(array)
	fmt.Println(array)
}
