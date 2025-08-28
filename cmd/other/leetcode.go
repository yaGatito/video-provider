package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func main1() {
	//area := areaOfMaxDiagonal([][]int{{6, 5}, {8, 6}, {2, 10}, {8, 1}, {9, 2}, {3, 5}, {3, 5}})
	diagonalOrder := findDiagonalOrder_TEST_NEW_METHOD([][]int{
		{1, 3, 7, 6, 8, 6, 2, 4},
		{8, 6, 4, 5, 3, 9, 3, 6}})

	//{1, 2, 3},
	//{4, 5, 6},
	//{7, 8, 9},
	//{7, 8, 9},
	//{5, 4, 1}})

	//{1, 2, 3, 3},
	//{4, 5, 6, 5},
	//{7, 8, 9, 2},
	//{7, 8, 9, 3},
	//{5, 4, 1, 7}})

	//{1, 2, 3, 5, 3},
	//{4, 5, 6, 8, 5},
	//{7, 8, 9, 3, 8},
	//{7, 8, 9, 4, 4},
	//{5, 4, 1, 2, 4}})

	//{1, 2, 3, 5, 3, 5, 3},
	//{4, 5, 6, 8, 5, 6, 5},
	//{7, 3, 7, 2, 9, 3, 8},
	//{7, 8, 9, 5, 1, 4, 2},
	//{5, 4, 1, 2, 4, 8, 5}})
	fmt.Println(diagonalOrder)
}

func generateSquare(sideLen int) [][]int {
	if sideLen < 1 {
		sideLen = rand.Intn(1<<8) + 2
	}
	out := make([][]int, rand.Intn(300))

	for i := range out {
		out[i] = make([]int, sideLen)
		for j := range out[i] {
			out[i][j] = rand.Intn(8) + 1
		}
	}

	return out
}

// 498. Diagonal Traverse
// https://leetcode.com/problems/diagonal-traverse/description/?envType=daily-question&envId=2025-08-25
// Given an m x n matrix mat, return an array of all the elements of the array in a diagonal order.

//Input: mat = [[1,2,3],[4,5,6],[7,8,9]]
//Output: [1,2,4,7,5,3,6,8,9]
//    j j j j j
// i [1,2,3]
// i [0,4,5,6]
// i [0,0,7,8,9]
// where index is shifting coefficient

func findDiagonalOrder(mat [][]int) []int {
	width := len(mat[0])
	height := len(mat)

	//for i := range mat {
	//	fmt.Println(mat[i])
	//}

	if height == 0 {
		return []int{}
	}

	if height == 1 {
		return mat[0]
	}
	wholeLen := height * width
	out := make([]int, wholeLen)
	if width == 1 {
		for i := range mat {
			out[i] = mat[i][0]
		}
	}

	isOrderPositive := false
	m := 0
	for j := 0; j < width*2-1; j++ {
		if j > 1 {
			isOrderPositive = !isOrderPositive
		}

		for i := 0; i <= j; i++ {
			if isOrderPositive {
				if j >= width {
					if i+1 >= width-(j-width) {
						break
					}
					out[m] = mat[(width - 1 - i)][i+(j-width)+1]
				} else {
					out[m] = mat[j-i][i]
				}
			} else {
				if j >= width {
					if i+1 >= width-(j-width) {
						break
					}
					out[m] = mat[i+(j-width)+1][(width - 1 - i)]
				} else {
					out[m] = mat[i][j-i]
				}
			}
			//fmt.Print(out[m])
			m++
		}
		//fmt.Print("\n")
	}

	return out
}

func findDiagonalOrder_TEST_NEW_METHOD(mat [][]int) []int {
	width := len(mat[0])
	height := len(mat)

	//for i := range mat {
	//	fmt.Println(mat[i])
	//}

	wholeLen := height * width
	out := make([]int, wholeLen)

	isOrderPositive := true
	m := 0
	iShift := 0
	for j := 0; j < height+width-1; j++ {
		if j > 1 {
			isOrderPositive = !isOrderPositive
		}

		if j >= width {
			iShift++
		}
		// positive width handler
		for i := iShift; i <= j; i++ {
			if isOrderPositive {
				if i >= height {
					continue
				}
				out[m] = mat[i][j-i]
			} else {
				if j-i+iShift >= height {
					continue
				}
				if i-iShift >= width {
					continue
				}
				out[m] = mat[j-i+iShift][i-iShift]
			}
			//fmt.Print(out[m])
			m++
		}
		//fmt.Print("\n")
	}

	return out
}

// in j on each turn should be the last element

// 3000. Maximum Area of Longest Diagonal Rectangle.
// https://leetcode.com/problems/maximum-area-of-longest-diagonal-rectangle/submissions/1748900131/?envType=daily-question&envId=2025-08-26
//Example 1:
//
//Input: dimensions = [[9,3],[8,6]]
//Output: 48
//Explanation:
//For index = 0, length = 9 and width = 3. Diagonal length = sqrt(9 * 9 + 3 * 3) = sqrt(90) ≈ 9.487.
//For index = 1, length = 8 and width = 6. Diagonal length = sqrt(8 * 8 + 6 * 6) = sqrt(100) = 10.
//So, the rectangle at index 1 has a greater diagonal length therefore we return area = 8 * 6 = 48.
//
//Example 2:
//
//Input: dimensions = [[3,4],[4,3]]
//Output: 12
//Explanation: Length of diagonal is the same for both which is 5, so maximum area = 12.

func areaOfMaxDiagonal(dimensions [][]int) int {
	maxArea := 0
	maxDiagonal := 0.0
	for _, d := range dimensions {
		diagonal := math.Sqrt((float64)(d[0]*d[0] + d[1]*d[1]))
		if diagonal > maxDiagonal {
			maxDiagonal = diagonal
			maxArea = d[0] * d[1]
		} else if diagonal == maxDiagonal {
			area := d[0] * d[1]
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

func curlyBrace(str string) string {
	out := strings.Replace(str, "[", "{", 1<<31-1)
	return strings.Replace(out, "]", "}", 1<<31-1)
}
