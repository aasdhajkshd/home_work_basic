package chessboard

import (
	"testing"
)

func TestDrawChessBoard(t *testing.T) {
	cases := []struct {
		axisX, axisY int    // axis: x, y
		sum          string // checkSum in MD5
	}{
		{8, 8, "d55b99581da7a69e9de8bdeea7f4c6bc"},
	}

	const (
		success = "\u2713"
		failed  = "\u2717"
	)

	t.Log("Тест для шахматной доски...")
	for n, c := range cases {
		checkSum := DrawChessBoard(c.axisX, c.axisY)
		if c.sum == checkSum {
			t.Logf("Тест #%d - \t%v", n+1, success)
		} else {
			t.Errorf("Тест #%d - \t%v", n+1, failed)
		}
	}
}
