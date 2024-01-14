package comparator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCompareBooks(t *testing.T) {
	t.Log("Устанавливаем значения для книг...")
	bookOne := Book{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden"}
	bookTwo := Book{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz"}

	const (
		success = "\u2713"
		failed  = "\u2717"
	)
	t.Log("Тест сравнения для книг\n")
	bookOneValues := reflect.ValueOf(bookOne)
	bookTwoValues := reflect.ValueOf(bookTwo)
	compareType := []CompareType{Id, Year, Size, Rate}

	for i := range compareType {
		var result, expectedResult bool
		result = CompareBooks(bookOne, bookTwo, i)
		switch bookOneValues.Field(i).Kind() { //nolint:exhaustive
		case reflect.Uint64:
			expectedResult = (bookOneValues.Field(i).Uint() > bookTwoValues.Field(i).Uint())
		case reflect.Float32:
			expectedResult = (bookOneValues.Field(i).Float() > bookTwoValues.Field(i).Float())
		case reflect.Int32:
			expectedResult = (bookOneValues.Field(i).Int() > bookTwoValues.Field(i).Int())
		default:
			fmt.Println("Тип для сравнения исключен из тестов")
		}
		if result == expectedResult {
			t.Logf("Тест #%d - \t%v", i+1, success)
		} else {
			t.Errorf("Тест #%d - \t%v (Ожидалось: %t)", i+1, failed, expectedResult)
		}
	}
}
