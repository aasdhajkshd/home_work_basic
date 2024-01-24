package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	cases := []struct {
		testSlice      []int
		testValue      int
		expectedResult int
		message        string
	}{
		// тест 1
		{
			testSlice:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			testValue:      5,
			expectedResult: 4,
			message:        "по-умолчанию",
		},
		// тест 2
		{
			testSlice:      []int{},
			testValue:      5,
			expectedResult: -1,
			message:        "пустой массив",
		},
		// тест 3, аналогичен тесту 6
		{
			testSlice:      []int{1},
			testValue:      5,
			expectedResult: -1,
			message:        "один элемент в массиве без совпадения",
		},
		// тест 4
		{
			testSlice:      []int{5},
			testValue:      5,
			expectedResult: 0,
			message:        "один элемент в массиве с совпадением",
		},
		// тест 5
		{
			testSlice:      []int{1, 2, 3, 4, 6, 7, 8, 9, 10},
			testValue:      5,
			expectedResult: -1,
			message:        "запрашиваемого значения нет в диапазоне",
		},
		// тест 6
		{
			testSlice:      []int{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			testValue:      5,
			expectedResult: -1,
			message:        "запрашиваемое значение вне диапазона",
		},
		// тест 7
		{
			testSlice:      []int{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			testValue:      -5,
			expectedResult: -1,
			message:        "отрицательное запрашиваемое значение вне диапазона",
		},
		// тест 8
		{
			testSlice:      []int{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			testValue:      -4,
			expectedResult: 0,
			message:        "искомое значение в начале массва",
		},
		// тест 9
		{
			testSlice:      []int{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			testValue:      3,
			expectedResult: 7,
			message:        "искомое значение последнее в массиве",
		},
	}
	for i, j := range cases {
		result := binarySearch(j.testSlice, j.testValue)
		assert.Equal(t, j.expectedResult, result, "Тест %d: %s", i+1, j.message)
	}
}
