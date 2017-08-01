package main

import (
	"fmt"
	"time"
)

type Point struct {
	x int
	y int
}

var canvas = [][]int{
	{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1},
}

var zeros = [][]int{}
var zeroMap = map[int]map[int]bool{}

func createRandCube() {

}

func main() {
	outshow(canvas)
	var start = time.Now().UnixNano()
	findZeroPoint()
	// fmt.Println(zeros)
	findMaxBlock()
	fmt.Println((time.Now().UnixNano() - start) / 1e3)
	outshow(canvas)
	/* var row = len(canvas)
	fmt.Println(row)
	outshow(canvas)
	var zeros = [][]int{} // x, y
	var KeyPoint = map[*Point]bool{}
	var start = time.Now().UnixNano()
	// 找出所有的零点
	for x, l := range canvas {
		var first = &Point{x, 0}
		KeyPoint[first] = true
		var zeroline = []int{}
		for y, p := range l {
			if p == 0 {
				zeroline = append(zeroline, y)
			}
		}
		zeros = append(zeros, zeroline)
	}
	fmt.Println("")
	// fmt.Println(zeros)
	// 开始找最大块
	for x := 0; x < len(canvas); x++ {
		var l = canvas[x]
		for y, p := range l {
			// fmt.Println(x, y, p)

			if p == 0 {
				zeros[x] = zeros[x][1:]
				// fmt.Println(x, y, zeros[x])
				continue
			}
			var max_in_row int
			// fmt.Println(zeros[x])
			if len(zeros[x]) > 0 {
				max_in_row = zeros[x][0] - y
			}
			if max_in_row <= 0 {
				max_in_row = len(l) - y
			}
			var max_in_col int = findColMax(x, y, zeros, row)

			// for k, l := range zeros[x+1:] {
			// 	if len(l) > 0 {
			// 		max_in_col = k - x
			// 	} else {
			// 		max_in_col = len(canvas) - x
			// 	}
			// }
			// fmt.Println(x, y, max_in_row, max_in_col)
			var min = min(max_in_row, max_in_col)
			// fmt.Println(x, y, zeros[x], max_in_row, max_in_col, min)
			if min > 1 {
				addOne(canvas, x, y, min)
			}
			// outshow(canvas)
		}
	}

	fmt.Println((time.Now().UnixNano() - start) / 1e3)
	outshow(canvas) */
}

func findMaxBlock() {
	var count = len(canvas)
	var blockMax = 1
	// 开始从0,0开始循环
	for x := 0; x < count-1; x++ {
		// 如果最大的已经超过离底线的距离 则不继续计算
		if blockMax > count-x {
			break
		}
		var row = canvas[x]
		var zeroI = 0
		// 循环一行
		for y := 0; y < len(row)-1; y++ {
			var p = row[y]
			var max int
			// fmt.Println(x)
			if len(zeros[x]) > zeroI {
				// fmt.Println(zeros[x], zeroI)
				max = zeros[x][zeroI] - y
				zeroI++
			} else {
				max = len(row) - y
			}

			if max > count-x {
				max = count - x
			}

			if max <= p {
				continue
			}
			var colmax, newY = findColMax2(x, y, max)
			max = min(colmax, max)
			// for i = 0; i < max; i++ {
			// if x+i == count-1 || y+i == len(row)-1 ||
			// 	checkZero(x+i+1, y+i+1) ||
			// 	checkZero(x+i, y+i+1) ||
			// 	checkZero(x+i+1, y+i) {
			// 	break
			// }
			// }
			// max = i
			if max > 1 {
				if max > blockMax {
					blockMax = max
				}
				flushPointMaxNumber(x, y, max)
			}
			if newY < count-1 {
				x = newY
			}
		}
	}
}

// 找出竖向最大 return max, 0y-index
func findColMax2(x, y, max int) (int, int) {
	// from next col to max col
	// 1	 1 1 0max
	// 0 1 1 <- done
	// 1 0 1 <- done
	// 1 1 0 <- done
	// ..... <- no

	// 1. in line < max
	// 		1.  zero.y > y
	//			zero.y - y > max return max
	//			zero.y - y < max return zero.y - y
	// 2. in line > max return max
	var colmax = 1
	var zeroy = y + max
	for _, z := range zeros[x+1 : x+max-1] {
		var done = false
		if len(z) > 0 {
			for _, j := range z {
				// 当找到y后面的0
				if j >= y {
					// 如果在max位置之后
					if j-y > max {
						colmax = max
						zeroy = y + max // 定位到该点max的0位置
						// return max, y + max
					} else {
						// 最大的为 0点到y的距离
						colmax = j - y
						if j-y < max/2 {
							// 如果0->y距离在max一半以内 定位到0点之后
							zeroy = j
						} else {
							zeroy = y + max // 定位到0max后
						}
						// return j - y, j
					}
					done = true
					break
				}
			}
		} else {
			colmax++
		}
		if done {
			break
		}
	}
	return colmax, zeroy
}

func flushPointMaxNumber(x, y, max int) {
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			if max > canvas[x+i][y+j] {
				canvas[x+i][y+j] = max
			}
		}
	}
}

func checkZero(x, y int) bool {
	if z, ok := zeroMap[x]; ok {
		if _, ok := z[y]; ok {
			return true
		}
	}
	return false
}

func findZeroPoint() {
	for x, l := range canvas {
		var zero = []int{}
		for y, p := range l {
			if p == 0 {
				zero = append(zero, y)
				if _, ok := zeroMap[x]; !ok {
					zeroMap[x] = map[int]bool{}
				}
				zeroMap[x][y] = true
			}
		}
		zeros = append(zeros, zero)
	}
	// fmt.Println(zeroMap)
}

func findColMax(x, y int, zeros [][]int, border int) int {
	for index := x + 1; index < len(zeros); index++ {
		if len(zeros[index]) > 0 {
			for _, k := range zeros[index] {
				if k >= y {
					return index - x
				}
			}
		}
	}
	return border - x
}

func addOne(canvas [][]int, x, y, max int) {
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			// fmt.Println("addto:", x+i, y+j, max)
			if max > canvas[x+i][y+j] {
				canvas[x+i][y+j] = max
			}
		}
	}
}

func outshow(vec [][]int) {
	for _, l := range vec {
		fmt.Println(l)
	}
}

func min(x, y int) int {
	if y < x {
		return y
	} else {
		return x
	}
}
