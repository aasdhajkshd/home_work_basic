package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			handler := http.HandlerFunc(HandlerHello)
			handler.ServeHTTP(w, req)
			assert.Equalf(t, http.StatusOK, w.Code, "Expected status code don't match")
			assert.Equalf(t, j.expectedResult, w.Body.String(), "Expected response result don't match")
		})
	}
}
