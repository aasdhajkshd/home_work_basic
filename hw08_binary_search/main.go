package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Range struct {
	min, max, num int
	nums          []int
}

//nolint:gosec
func (r *Range) Random() ([]int, error) {
	if r.max <= 0 || r.max-1 == r.min {
		return nil, fmt.Errorf("здесь 100%% результат") // плохой диапазон будет, просто отдадим ошибку
	}
	nums := make(map[int]struct{}, r.num)
	for i := 0; i < r.num; i++ {
		v := r.min + rand.Intn(r.max-r.min+1)
		if _, found := nums[v]; !found {
			nums[v] = struct{}{}
			r.nums = append(r.nums, v)
		}
	}
	// fmt.Printf("Random r.nums: %+v\n", r.nums)
	return r.nums, nil
}

// пакет slices появился Go 1.21
func contains(s []int, v int) bool {
	for _, j := range s {
		if j == v {
			return true
		}
	}
	return false
}

// на вход принимает слайс чисел и число для поиска.
func binarySearch(inputSlice []int, searchValue int) int {
	lenSlice := len(inputSlice)
	if lenSlice == 0 { // нулевой массив
		return -1
	}
	lowValue, highValue := 0, lenSlice
	for lowValue <= highValue { // цикл-условие аналог "while"
		j := (lowValue + highValue) / 2
		if j == lenSlice { // запрашиваемое значение вне допустимого диапазона
			break
		}
		if inputSlice[j] == searchValue { // нужно выше, если прямое попадание посередине
			return j
		}
		if searchValue > inputSlice[j] {
			lowValue = j + 1
		} else {
			highValue = j - 1
		}
		// fmt.Println("Индекс:", j, "слева:", lowValue, "справа:", highValue, "значение:", inputSlice[j], "искомое:", searchValue) //nolint:lll
	}
	return -1
}

func main() {
	var r Range
	r.min = 0
	r.max = 100
	r.num = 40
	var inputSlice []int
	var searchValue int
	// fmt.Println(bs.binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5))
	// os.Exit(0)

	if v, err := r.Random(); err != nil {
		inputSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		searchValue = 11
	} else {
		inputSlice = v
		sort.Ints(inputSlice) // slices.Sort(inputSlice)
	}
	fmt.Println("Поисковая комбинация:", inputSlice)
	// fmt.Println("Main:", r)
	fmt.Printf("Укажите цифру для поиска в диапазоне от %d до %d\n> ", r.min, r.max)
	fmt.Scanln(&searchValue)
	if searchValue >= r.min && searchValue <= r.max && contains(inputSlice, searchValue) { // slices.Contains(inputSlice, searchValue)
		indexSlice := binarySearch(inputSlice, searchValue)
		if indexSlice > 0 {
			fmt.Printf("Индекс: %d, значение: %d\n", indexSlice, inputSlice[indexSlice])
		} else {
			fmt.Printf("Индекс: %d\n", indexSlice)
		}
	} else {
		fmt.Printf("Значение не найдено в массиве: %d\n", searchValue)
	}
}
