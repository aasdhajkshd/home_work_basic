package chessboard

import (
	"fmt"
	"testing"
)

func TestDrawChessBoard(t *testing.T) {
	cases := []struct {
		axisX, axisY int    // axis: x, y
		expectedResult string
		testDescription  string
	}{
		{8, 8, "|#| |#| |#| |#| |\n| |#| |#| |#| |#|\n|#| |#| |#| |#| |\n| |#| |#| |#| |#|\n|#| |#| |#| |#| |\n| |#| |#| |#| |#|\n|#| |#| |#| |#| |\n| |#| |#| |#| |#|\n", "доска 8x8"},
		{5, 5, "|#| |#| |#|\n| |#| |#| |\n|#| |#| |#|\n| |#| |#| |\n|#| |#| |#|\n", "доска 5x5"}, // это пример от обратного тестирования 5x8
	}

	const (
		success = "\u2713"
		failed  = "\u2717"
	)

	t.Log("Тест для шахматной доски...")
	for n, c := range cases {
		t.Run(c.testDescription, func(t *testing.T) {
			result := DrawChessBoard(c.axisX, c.axisY)
			fmt.Println(result)
			fmt.Println(c.expectedResult)
			if c.expectedResult == result {
				t.Logf("Тест #%d - \t%v", n+1, success)
			} else {
				t.Errorf("Тест #%d - \t%v", n+1, failed)
			}
		})
	}
}
