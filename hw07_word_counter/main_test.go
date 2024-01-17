package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleText string = `Вот дом,
Который построил Джек.
А это пшеница,
Которая в темном чулане хранится
В доме,
Который построил Джек.
А это веселая птица-синица,
Которая часто ворует пшеницу,
Которая в темном чулане хранится
В доме,
Который построил Джек.`

func TestCountWords(t *testing.T) {
	checkText := map[string]int{
		"а":            2,
		"в":            4,
		"веселая":      1,
		"ворует":       1,
		"вот":          1,
		"джек":         3,
		"дом":          1,
		"доме":         2,
		"которая":      3,
		"который":      3,
		"построил":     3,
		"птица-синица": 1,
		"пшеница":      1,
		"пшеницу":      1,
		"темном":       2,
		"хранится":     2,
		"часто":        1,
		"чулане":       2,
		"это":          2,
	}
	testText, err := countWords(exampleText)

	if err == nil {
		assert.Equal(t, testText, checkText, "Ошибка проверки:\n%v\n%d", testText, checkText)
	} else {
		assert.Error(t, err)
	}
}
