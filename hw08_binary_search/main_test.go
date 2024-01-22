package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimSpaceCR(t *testing.T) {
	result := string(trimSpaceCR([]byte(" Привет Отус! \r\n")))
	expectedResult := "Привет Отус!"
	assert.Equal(t, expectedResult, result, "Ошибка удаления каретки и пробелов")
}

func TestScanComma(t *testing.T) {
	testInput := []string{
		"Санкт-Петербург, Москва, Ставрополь",
		"Старый Оскол",
	}
	for _, j := range testInput {
		s := bufio.NewScanner(strings.NewReader(j))
		s.Split(scanComma)
		result := bytes.Split([]byte(j), []byte(", ")) // с учётом допущений двух знаков
		expectedResult := [][]byte{}
		for s.Scan() {
			expectedResult = append(expectedResult, s.Bytes())
		}
		// t.Logf("\"%s\": %q\n", j, expectedResult)
		// t.Logf("\"%s\": %q\n", j, result)
		assert.Equal(t, expectedResult, result, "проверка на разделение по запятым провалилась")
	}
}

func TestFindMatch(t *testing.T) {
	testInput := []string{
		"Архангельск",
		"Санкт-Петербург",
		"Сергиев Посад",
	}
	byteWords, err := readFile("list_test.enc")
	if err != nil {
		t.Fatal("ох, ошибка чтения файла... ")
	}
	for _, j := range testInput {
		if !findMatch(j, byteWords) {
			t.Fail()
		}
	}
}
