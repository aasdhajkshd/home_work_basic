package chessboard

import (
	"testing"
)

func TestDrawChessBoard(t *testing.T) {
	cases := []struct {
		axisX, axisY int  // axis: x, y
		result       bool // false - expected fail , true - ok
	}{
		{0, 4, false}, // here we can set to false to "test" this test
		{-1, 3, false},
		{14, 5, false},
		{8, 8, true},
	}

	const (
		success = "\u2713"
		failed  = "\u2717"
	)

	t.Log("Здесь проводим тестирование значений, " +
		"которые должны быть в разрешенном диапазоне от 2 до 12.\n")
	for i, c := range cases {
		result := DrawChessBoard(c.axisX, c.axisY)
		if (result != nil) == c.result {
			t.Errorf("Тест #%d - \t%v (Ожидалось неудача: %t)", i+1, failed, c.result)
		} else {
			t.Logf("Тест #%d - \t%v", i+1, success)
		}
	}
}

func BenchmarkDrawChessBoard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := DrawChessBoard(8, 8); err != nil {
			b.Log(err)
		}
	}
}
