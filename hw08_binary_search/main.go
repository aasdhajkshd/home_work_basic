package main

import (
	"fmt"
	"math/rand"
	"slices"
)

type Range struct {
	min, max, num int
	nums          []int
}

type BinaryStruct struct {
	lowValue, highValue int
}

//nolint:gosec
func (r *Range) Random() ([]int, error) {
	if r.max <= 0 || r.max-1 == r.min {
		return nil, fmt.Errorf("здесь 100%% результат") // плохой диапазон будет, просто отдадим ошибку
	}
	nums := make(map[int]bool, r.num)
	for i := 0; i < r.num; i++ {
		v := r.min + rand.Intn(r.max-r.min+1)
		if _, found := nums[v]; found {
			nums[v] = true // для понимания количества уникальных чисел
		} else {
			nums[v] = false
			// реализовывал и на crypto, для домашнего, без велосипеда, тут логига другая
			r.nums = append(r.nums, v)
		}
	}
	// fmt.Printf("Random r.nums: %+v\n", r.nums)
	return r.nums, nil
}

// на вход принимает слайс чисел и число для поиска.
func (b *BinaryStruct) binarySearch(inputSlice []int, searchValue int) int {
	lenSlice := len(inputSlice)
	if lenSlice == 0 { // нулевой массив
		return -1
	}
	b.lowValue, b.highValue = 0, lenSlice
	for b.lowValue <= b.highValue { // цикл-условие аналог "while"
		j := (b.lowValue + b.highValue) / 2
		if j == lenSlice { // запрашиваемое значение вне допустимого диапазона
			break
		}
		if inputSlice[j] == searchValue { // нужно выше, если прямое попадание посередине
			return j
		}
		if searchValue > inputSlice[j] {
			b.lowValue = j + 1
		} else {
			b.highValue = j - 1
		}
		fmt.Println("Индекс:", j, "слева:", b.lowValue, "справа:", b.highValue, "значение:", inputSlice[j], "искомое:", searchValue) //nolint:lll
	}
	return -1
}

func main() {
	var r Range
	r.min = 0
	r.max = 100
	r.num = 40
	bs := BinaryStruct{}
	var inputSlice []int
	var searchValue int
	// fmt.Println(bs.binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5))
	// os.Exit(0)
	if v, err := r.Random(); err != nil {
		inputSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		searchValue = 11
	} else {
		inputSlice = v
		slices.Sort(inputSlice)
	}
	fmt.Println("Поисковая комбинация:", inputSlice)
	// fmt.Println("Main:", r)
	fmt.Printf("Укажите цифру для поиска в диапазоне от %d до %d\n> ", r.min, r.max)
	fmt.Scanln(&searchValue)
	if searchValue >= r.min && searchValue <= r.max {
		indexSlice := bs.binarySearch(inputSlice, searchValue)
		if indexSlice < 0 {
			fmt.Printf("Значение не найдено: %d\n", searchValue)
		} else {
			fmt.Printf("Индекс: %d, значение: %d\n", indexSlice, inputSlice[indexSlice])
		}
	} else {
		fmt.Println("указано значение вне допустимого диапазона")
	}
}
