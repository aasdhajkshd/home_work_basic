package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	jsonFile = "data/data_test.json"
)

func TestReadJSON(t *testing.T) {
	expectedResult := Data{
		Users: []User{
			{ID: 1, Name: "Vasiliy Ivanov", Email: "vi@mail.ru", Password: ""},
		},
	}

	data := ReadJSON(jsonFile)
	assert.Equal(t, data.Users, expectedResult.Users)
}
