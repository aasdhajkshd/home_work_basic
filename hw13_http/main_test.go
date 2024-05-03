package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	_ "time"

	server "github.com/aasdhajkshd/home_work_basic/hw13_http/server"
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

func TestRunServer(t *testing.T) {
	expectedResult := "ehlo test"
	t.Run("Test POST handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/hello?name=test", nil)
		w := httptest.NewRecorder()
		server.HandlerHello(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if string(data) != expectedResult {
			t.Errorf("Expected %s but got %v", expectedResult, string(data))
		}
	})
}
