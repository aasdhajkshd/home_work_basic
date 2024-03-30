package main

import (
	_ "fmt"
	"testing"
	"time"

	book "github.com/aasdhajkshd/home_work_basic/hw09_serialize/book"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

type TestBook struct {
	TestName       string
	TestBook       Book
	ExpectedResult string
}

func TestMarshallerJSON(t *testing.T) {
	cases := []TestBook{
		{
			TestName: "Test 1",
			TestBook: Book{
				Id:      1234567890,
				Year:    2024,
				Size:    101,
				Rate:    4.1,
				Title:   "Title",
				Author:  "Author",
				Updated: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC), //nolint:gofumpt
			},
			ExpectedResult: "{\"updated\":1546304400,\"id\":1234567890,\"year\":2024,\"size\":101,\"rate\":4.1,\"title\":\"Title\",\"author\":\"Author\"}", //nolint:lll
		},
	}
	for _, j := range cases {
		t.Run(j.TestName, func(t *testing.T) {
			jsonResult, err := j.TestBook.MarshalJSON()
			if err != nil {
				t.Errorf("Error marshalling book via json: %v", err)
				return
			}
			j.TestBook.UnmarshalJSON(jsonResult)
			assert.JSONEqf(t, j.ExpectedResult, string(jsonResult), "Expected JSON and actual results don't match")
		})
	}
}

func TestMarshallerProto(t *testing.T) {
	cases := &struct {
		TestName       string
		TestMsg        book.Message
		ExpectedResult string
	}{
		TestName: "Test 1",
		TestMsg: book.Message{
			Id:     1234567890,
			Year:   2024,
			Size:   101,
			Rate:   4.1,
			Title:  "Title",
			Author: "Author",
		},
		ExpectedResult: "id:1234567890 year:2024 size:101 rate:4.1 title:\"Title\" author:\"Author\"",
	}
	t.Run(cases.TestName, func(t *testing.T) {
		protoResult, err := proto.Marshal(&cases.TestMsg)
		if err != nil {
			t.Errorf("Error marshalling book via protobuf: %v", err)
			return
		}
		resultMessage := &book.Message{}
		proto.Unmarshal(protoResult, resultMessage)
		assert.Equal(t, cases.ExpectedResult, resultMessage.String(), "Expected and actual results don't match")
	})
}
