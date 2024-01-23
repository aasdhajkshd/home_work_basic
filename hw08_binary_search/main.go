package main

import (
	"fmt"
)

type binaryStruct struct {
	lowValue, highValue int
}

func (b *binaryStruct) binarySearch(inputSlice []int, searchValue int) int {
	b.lowValue = 0
	b.highValue = len(inputSlice)
	for i := b.lowValue; i < b.highValue; i++ {
		j := (b.lowValue + b.highValue) / 2
		if inputSlice[j] == searchValue {
			return j
		}
		if searchValue > inputSlice[j] {
			b.lowValue = j + 1
		} else {
			b.highValue = j - 1
		}
	}
	return -1
}

func main() {
	bs := binaryStruct{}
	inputSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	searchValue := 11
	fmt.Println(bs.binarySearch(inputSlice, searchValue)) // -1, если числа нет.
}
