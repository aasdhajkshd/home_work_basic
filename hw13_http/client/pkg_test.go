package client

import (
	"testing"
)

func TestRequestURL(t *testing.T) {
	testCases := []struct {
		name, url, method string
		expectedResult    bool
	}{
		{
			name:           "URL is available",
			url:            "https://go.dev/play",
			method:         "GET",
			expectedResult: false,
		},
		{
			name:           "URL localhost is refused",
			url:            "http://localhost:8081/hello",
			method:         "GET",
			expectedResult: true,
		},
	}
	for _, j := range testCases {
		t.Run(j.name, func(t *testing.T) {
			resp, err := RequestURL(j.method, j.url, "plain/text", 10)
			if (err != nil) != j.expectedResult {
				t.FailNow()
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close() // Закрываем тело ответа, требуется response body must be closed
			}
		})
	}
}
