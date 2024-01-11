package comparator

import (
	"testing"
)

func TestCompareBooks(t *testing.T) {
	t.Log("Устанавливаем значения для книг...")
	books := []struct {
		id            uint64
		year, size    rune
		rate          float32
		title, author string
	}{
		{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden"}, // here we can set to false to "test" this test
		{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz"},
	}

	const (
		success = "\u2713"
		failed  = "\u2717"
	)

	t.Log("Здесь проводим тестирование книг для сравнения\n")
	for i := uint8(1); i < 4; i++ {
		r := CompareBooks(books[0], books[1], i)
		if r != nil {
			t.Errorf("Тест сравнения -\t%v", failed)
		} else {
			t.Logf("Тест сравнения -\t%v", success)
		}
	}
}
