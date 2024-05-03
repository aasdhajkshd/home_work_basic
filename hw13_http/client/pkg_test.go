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
			name:           "S is unavailable", // блокировка на уровне провайдера
			url:            "https://www.digitalocean.com",
			method:         "GET",
			expectedResult: true,
		},
		{
			name:           "URL is available",
			url:            "https://phet-dev.colorado.edu/html/build-an-atom/0.0.0-3/simple-text-only-test-page.html",
			method:         "GET",
			expectedResult: false,
		},
		{
			name:           "localhost URL is refused",
			url:            "https://localhost:8083",
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
