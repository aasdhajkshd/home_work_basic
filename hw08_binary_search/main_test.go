package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	testSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	testValue := 5      // поставить 11 для получения ошибки
	expectedResult := 4 // -1 в слайсе
	result := binarySearch(testSlice, testValue)
	assert.Equal(t, expectedResult, result, "Возникла ошибка проверки функции двоичного поиска")
}
