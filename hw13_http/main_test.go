package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	client "github.com/aasdhajkshd/home_work_basic/hw13_http/client"
	server "github.com/aasdhajkshd/home_work_basic/hw13_http/server"
	"github.com/stretchr/testify/assert"
)

func TestRequestURL(t *testing.T) {
	testCases := []struct {
		name, url, method string
		expectedResult    bool
	}{
		{
			name:           "URL is unavailable", // блокировка на уровне провайдера
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
	}
	for _, j := range testCases {
		t.Run(j.name, func(t *testing.T) {
			resp, err := client.RequestURL(j.method, j.url, "plain/text", 10)
			if (err != nil) != j.expectedResult {
				t.FailNow()
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close() // Закрываем тело ответа, требуется response body must be closed
			}
		})
	}
}

func TestHandler(t *testing.T) {
	testCases := []struct {
		name, path, method string
		expectedResult     string
	}{
		{
			name:           "Server handler test",
			path:           "/hello",
			method:         "GET",
			expectedResult: "Привет, мой друг!",
		},
		{
			name:           "Server handler test",
			path:           "/hello",
			method:         "POST",
			expectedResult: "Привет, пишете письмо?",
		},
	}
	for _, j := range testCases {
		t.Run(j.name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, j.method, j.path, nil)
			if err != nil {
				t.Fatalf("Failed to create %s request: %v", j.method, err)
			}

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(server.HandlerHello)
			handler.ServeHTTP(w, req)
			assert.Equalf(t, http.StatusOK, w.Code, "Expected status code don't match")
			assert.Equalf(t, j.expectedResult, w.Body.String(), "Expected response result don't match")
		})
	}
}
