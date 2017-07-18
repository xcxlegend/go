package main

import (
	"fmt"
)

/**
 * 堆排序
 * @param {[type]} array []int [description]
 */
func HeapSort(array []int) {

	var count = len(array)

	for i := count/2 - 1; i >= 0; i-- {
		heap_ajust(array, i, count-1)
	}

	for i := count - 1; i >= 0; i-- {
		array[0], array[i] = array[i], array[0]
		heap_ajust(array, 0, i-1)
	}

}

/**
 * 重建最大堆
 * @param  {[type]} array []int         [description]
 * @param  {[type]} start [description]
 * @param  {[type]} end   int           [description]
 * @return {[type]}       [description]
 */
func heap_ajust(array []int, start, end int) {

	var temp = array[start]

	// ! 重点  每次循环对子节点的重排
	for i := 2*(start+1) - 1; i <= end; i = 2*(i+1) - 1 {
		if i < end && array[i] < array[i+1] {
			i++
		}
		if array[i] <= temp {
			break
		}
		array[start] = array[i]
		start = i
	}
	array[start] = temp

}

func main() {
	var array = []int{1, 99, 23, 75, 65, 32, 12, 35, 77, 4}
	HeapSort(array)
	fmt.Println(array)
}
