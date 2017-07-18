package main

import (
	"fmt"
)

/**
 * 希尔排序
 * ! 重点在于选择一个间隔数
 * @param {[type]} array []int [description]
 */
func HillSort(array []int) {

	var count = len(array)

	// 选择每次一半
	var gap = count / 2

	for gap > 0 {

		for i := gap; i < count; i++ {

			var temp = i
			for temp-gap >= 0 {

				if array[temp] > array[temp-gap] {
					break
				}
				array[temp], array[temp-gap] = array[temp-gap], array[temp]

				temp -= gap
			}

		}

		gap /= 2
	}

}

func main() {
	var array = []int{1, 99, 23, 75, 65, 32, 12, 35, 77, 4}
	HillSort(array)
	fmt.Println(array)
}
