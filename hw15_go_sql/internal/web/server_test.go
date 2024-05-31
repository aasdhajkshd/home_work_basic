package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/database"
	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	testDsn string = "postgres://manager:password@postgres:5432/onlinestore_test?sslmode=disable&pool_max_conns=100"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "POST user",
			method:         http.MethodPost,
			path:           "/user",
			body:           models.User{Name: "John Doe", Email: "john@example.com", Password: "12345678"},
			expectedStatus: http.StatusCreated,
			expectedResult: `{"id":1}`,
		},
		{
			name:           "POST product 1",
			method:         http.MethodPost,
			path:           "/product",
			body:           models.Product{Name: "Phone", Price: 500.0},
			expectedStatus: http.StatusCreated,
			expectedResult: `{"id":1}`,
		},
		{
			name:           "POST product 2",
			method:         http.MethodPost,
			path:           "/product",
			body:           models.Product{Name: "Laptop", Price: 1500.0},
			expectedStatus: http.StatusCreated,
			expectedResult: `{"id":2}`,
		},
		{
			name:           "POST order",
			method:         http.MethodPost,
			path:           "/order",
			body:           models.Order{UserID: 1, OrderDate: time.Unix(0, 0), TotalAmount: 200, Products: models.Products{List: []models.Product{{ID: 1, Name: "Phone", Price: 500.0}, {ID: 2, Name: "Laptop", Price: 1500.0}}}}, //nolint:lll
			expectedStatus: http.StatusCreated,
			expectedResult: `{"id":1}`,
		},
		{
			name:           "GET order",
			method:         http.MethodGet,
			path:           "/order?id=1",
			expectedStatus: http.StatusOK,
			expectedResult: `{"id":1,"userID":0,"orderDate":"1970-01-01T03:00:00Z","totalAmount":200,"products":{"List":[{"id":1,"name":"Phone","price":500},{"id":2,"name":"Laptop","price":1500}]}}`, //nolint:lll
		},
	}

	db, err := database.NewPgxPool(testDsn)
	if err != nil {
		t.Fatalf("Failed to create DB connection: %v", err)
	}
	defer db.Close()
	db.InitDB()
	db.ResetDB()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var req *http.Request

			if tc.method == http.MethodGet {
				req, err = http.NewRequestWithContext(ctx, tc.method, tc.path, nil)
			} else {
				bodyBytes, _ := json.Marshal(tc.body)
				req, err = http.NewRequestWithContext(ctx, tc.method, tc.path, bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			}

			if err != nil {
				t.Fatalf("Failed to create %s request: %v", tc.method, err)
			}

			app := NewApplication(db)
			w := httptest.NewRecorder()
			switch tc.path {
			case "/user":
				app.handlerUser(w, req)
			case "/product":
				app.handlerProduct(w, req)
			default: // Здесь нужно Router().NotFoundHandler
				app.handlerOrder(w, req)
			}
			assert.Equalf(t, tc.expectedStatus, w.Code, "Expected status code don't match")
			assert.JSONEqf(t, tc.expectedResult, w.Body.String(), "Expected response result don't match")
		})
	}
	db.ResetDB()
}
