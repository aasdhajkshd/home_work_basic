package main

import (
	"fmt"
	"math/rand"
	"slices"
	_ "sort"
)

type rangeRand struct {
	min, max, num int
	nums          []int
}

//nolint:gosec
func (r *rangeRand) Random() ([]int, error) {
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
func binarySearch(r []int, i int) int {
	k := make([]int, len(r))
	copy(k, r) // да, нужно клонировать или новый append slice
	slices.Sort(r)
	midValue := len(r) / 2
	lowValue := slices.Min(r)
	highValue := slices.Max(r)
	// fmt.Println("binarySearch:", lowValue, midValue, highValue, r)
	for _, j := range r {
		if i > j {
			lowValue = midValue + 1
		} else {
			highValue = midValue - 1
		}
		midValue = (lowValue + highValue) // 2
		if j == i {
			fmt.Println("! по классике после сортировки нашли - ", i, "из", r)
		}
	}
	f := slices.Index(k, i)
	fmt.Println("! в пакете slices без сортировки позиция", f+1, "из", k)
	return f // возвращает этот метод позиция найденного числа
}

func main() {
	var r rangeRand
	r.min = 9
	r.max = 29
	r.num = 6
	var e []int
	if v, err := r.Random(); err != nil {
		e = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	} else {
		e = v
	}
	fmt.Println("поисковая комбинация:", e)
	// fmt.Println("Main:", r)
	var i int
	fmt.Printf("Прекрасная игра в угадай цифру, укажите её от %d до %d\n> ", r.min, r.max)
	fmt.Scanln(&i)
	if i >= r.min && i <= r.max {
		fmt.Println("вот что мы нашли...:", binarySearch(e, i)) // или -1, если числа нет.
	} else {
		fmt.Println("указано значение вне диапазона")
	}
}
