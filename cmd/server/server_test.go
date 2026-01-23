package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bidntb/metrics/internal/middleware"
	"bidntb/metrics/internal/service/handler"
	"bidntb/metrics/internal/service/metrics"
	"bidntb/metrics/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	storageInstance := storage.NewMemStorage()
	metricsSvc := metrics.NewService(storageInstance)
	h := handler.NewHandler(metricsSvc)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler(h.NotFoundHandler, h.BadRequestHandler))

	r.POST("/update/:type/:name/:value", h.UpdateMetric)
	r.POST("/update/", h.UpdateMetric)
	r.GET("/value/:type/:name", h.GetValue)
	r.GET("/", h.ListMetrics)

	r.POST("/update/counter", h.NotFoundHandler)
	r.POST("/update/gauge/", h.NotFoundHandler)
	r.POST("/update/gauge", h.NotFoundHandler)
	r.NoRoute(h.NotFoundHandler)

	return r
}

func TestServerRoutes(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name         string
		method       string
		path         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "List all metrics",
			method:       http.MethodGet,
			path:         "/",
			expectedCode: http.StatusOK,
		},
		{
			name:   "Update gauge JSON",
			method: http.MethodPost,
			path:   "/update/",
			body: `{
				"id": "test_gauge",
  				"type": "gauge",
  				"value": 123.45
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Update counter path params",
			method:       http.MethodPost,
			path:         "/update/counter/test_counter/5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get metric value path params",
			method:       http.MethodGet,
			path:         "/value/gauge/test_gauge",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid value - bad request",
			method:       http.MethodPost,
			path:         "/update/counter/test/none",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Not found route",
			method:       http.MethodPost,
			path:         "/unknown",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestUpdateMetricJSON(t *testing.T) {
	r := setupTestRouter()

	body := `{"id":"requests","type":"counter","delta":10}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCounterAccumulation(t *testing.T) {
	r := setupTestRouter()
	reqs := []string{
		`{"id":"requests","type":"counter","delta":10}`,
		`{"id":"requests","type":"counter","delta":20}`,
	}

	for _, body := range reqs {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/value/counter/requests", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "30", w.Body.String())
}

func TestValidationErrors(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name         string
		path         string
		body         string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Missing JSON",
			path:         "/update/",
			body:         `{"type":"gauge"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid counter value",
			path:         "/update/counter/test/none",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "invalid value for counter",
		},
		{
			name:         "Invalid gauge value",
			path:         "/update/gauge/test/none",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "invalid value for gauge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedMsg != "" {
				assert.Contains(t, w.Body.String(), tt.expectedMsg)
			}
		})
	}
}
