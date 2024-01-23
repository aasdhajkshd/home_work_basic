package main

import (
	"fmt"
)

type binaryStruct struct {
	lowValue, midValue, highValue int
}

func (b *binaryStruct) binarySearch(inputSlice []int, searchValue int) int {
	b.lowValue = 0
	b.highValue = len(inputSlice) - 1 // len возвращает количество
	for i, j := range inputSlice {
		if searchValue > j {
			b.lowValue = b.midValue + 1
		} else {
			b.highValue = b.midValue - 1
		}
		b.midValue = (b.lowValue + b.highValue) / 2
		if j == searchValue {
			return i
		}
	}
	return -1
}

func main() {
	bs := binaryStruct{}
	inputSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	searchValue := 5
	fmt.Println(bs.binarySearch(inputSlice, searchValue)) // -1, если числа нет.
}
